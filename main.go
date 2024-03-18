// Package classification VK-Intern 2024
//
// # Documentation for VK-Intern 2024 Ширшов Артём
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
// Produces:
// - application/json
// SecurityDefinitions:
//  key:
//   type: apiKey
//   in: header
//   name: Authorization
//
// swagger:meta
package main

import (
	"net/http"
	actorDelivery "vk-intern_test-case/internal/actor/delivery"
	actorRepository "vk-intern_test-case/internal/actor/repository"
	"vk-intern_test-case/internal/middleware"
	"vk-intern_test-case/utils/database"

	log "github.com/sirupsen/logrus"

	filmDelivery "vk-intern_test-case/internal/film/delivery"
	filmRepository "vk-intern_test-case/internal/film/repository"

	openApiMiddleware "github.com/go-openapi/runtime/middleware"
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

	authMw := middleware.NewAuthMiddleware(dbPool)

	r := http.NewServeMux()
	actorsHandler := http.HandlerFunc(aD.HandleActors)
	r.Handle("/actors", authMw.MiddlewareCheckAdmin(actorsHandler))
	r.Handle("/actors/", authMw.MiddlewareCheckAdmin(actorsHandler))

	filmsHandler := http.HandlerFunc(fD.HandleFilms)
	r.Handle("/films", authMw.MiddlewareCheckAdmin(filmsHandler))
	r.Handle("/films/", authMw.MiddlewareCheckAdmin(filmsHandler))

	r.HandleFunc("/film", fD.HandleFilm)

	opts := openApiMiddleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := openApiMiddleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	log.Info("Server started")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Error(err)
		return
	}
}
