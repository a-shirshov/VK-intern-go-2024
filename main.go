// Package classification Actor API
//
// # Documentation for Actor API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta
package main

import (
	"net/http"
	actorDelivery "vk-intern_test-case/internal/actor/delivery"
	actorRepository "vk-intern_test-case/internal/actor/repository"
	"vk-intern_test-case/utils/database"

	log "github.com/sirupsen/logrus"

	filmDelivery "vk-intern_test-case/internal/film/delivery"
	filmRepository "vk-intern_test-case/internal/film/repository"

	"github.com/go-openapi/runtime/middleware"
)

func main() {
	logLevel := log.DebugLevel
	log.SetLevel(logLevel)

	dbPool, err := database.InitPostgres()
	if err != nil {
		log.Error(err)
		return
	}
	defer dbPool.Close()

	fR := filmRepository.NewFilmRepository(dbPool)
	fD := filmDelivery.NewFilmDelivery(fR)

	aR := actorRepository.NewActorRepository(dbPool)
	aD := actorDelivery.NewActorDelivery(aR)

	r := http.NewServeMux()
	r.HandleFunc("/actors", aD.AddActor)
	r.HandleFunc("/actors/", aD.HandleActor)
	
	r.HandleFunc("/films", fD.AddFilm)
	r.HandleFunc("/films/", fD.HandleFilm)

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	log.Info("Server started")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Error(err)
		return
	}
}
