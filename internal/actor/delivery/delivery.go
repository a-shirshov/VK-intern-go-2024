package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"vk-intern_test-case/internal/actor"
	"vk-intern_test-case/models"
	"vk-intern_test-case/utils/response"

	log "github.com/sirupsen/logrus"
)

const logMessage = "actor:delivery:"

type actorDelivery struct {
	actorRepo actor.ActorRepository
}

func NewActorDelivery(aR actor.ActorRepository) *actorDelivery {
	return &actorDelivery{
		actorRepo: aR,
	}
}

func (aD *actorDelivery) HandleActors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		aD.GetActors(w, r)
	case http.MethodPost:
		aD.AddActor(w, r)
	case http.MethodPut:
		aD.UpdateActor(w, r)
	case http.MethodDelete:
		aD.DeleteActor(w, r)
	}
}

// swagger:route POST /actors Actors addActor
// Добавляет актёра в систему
// security:
// - key:
// responses:
//
//	200: actor
//	400: basicResponse
//  401: basicResponse
//	500: basicResponse
func (aD *actorDelivery) AddActor(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "AddActor:"
	log.Info(message + "started")
	jsonEnc := response.MakeJsonEncoder(w)
	var actor models.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}
	resultActor, err := aD.actorRepo.AddActor(&actor)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteResponse(w, jsonEnc, http.StatusOK, resultActor)
}

// swagger:route PUT /actors/{id} Actors updateActor
// Обновляет информацию об актёре. На вход полная информация.
// security:
// - key:
// responses:
//
//	200: basicResponse
//	400: basicResponse
//  401: basicResponse
//  500: basicResponse
func (aD *actorDelivery) UpdateActor(w http.ResponseWriter, r *http.Request) {
	jsonEnc := response.MakeJsonEncoder(w)
	id := strings.TrimPrefix(r.URL.Path, "/actors/")
	actorID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	var actor models.Actor
	err = json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	err = aD.actorRepo.UpdateActor(actorID, &actor)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteBasicResponse(w, jsonEnc, http.StatusOK, "OK")
}

// swagger:route DELETE /actors/{id} Actors deleteActor
// Удаляет актёра из системы.
// security:
// - key:
// responses:
//
//	200: basicResponse
//	400: basicResponse
//  401: basicResponse
//	500: basicResponse
func (aD *actorDelivery) DeleteActor(w http.ResponseWriter, r *http.Request) {
	jsonEnc := response.MakeJsonEncoder(w)

	id := strings.TrimPrefix(r.URL.Path, "/actors/")
	actorID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	err = aD.actorRepo.DeleteActor(actorID)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteBasicResponse(w, jsonEnc, http.StatusOK, "OK")
}

// swagger:route GET /actors Actors getActors
// Возращает список актёров с их фильмами.
// responses:
//
//	200: []actorWithFilms
//	500: basicResponse
func (aD *actorDelivery) GetActors(w http.ResponseWriter, r *http.Request) {
	jsonEnc := response.MakeJsonEncoder(w)

	resultActors, err := aD.actorRepo.GetActors()
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteResponse(w, jsonEnc, http.StatusOK, resultActors)
}
