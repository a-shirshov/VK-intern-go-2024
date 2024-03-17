package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"vk-intern_test-case/internal/actor"
	"vk-intern_test-case/models"

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

func (aD *actorDelivery) HandleActor(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		aD.AddActor(w, r)
	case http.MethodPut:
		aD.UpdateActor(w, r)
	case http.MethodDelete:
		aD.DeleteActor(w, r)
	}
}

// swagger:route POST /actors Actors addActor
// Returns a status of operation adding an actor
// responses:
//	200: basicResponse
//  400: basicResponse
func (aD *actorDelivery) AddActor(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "AddActor:"
	log.Info(message + "started")
	var actor models.Actor
	jsonEnc := json.NewEncoder(w)
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	resultActor, err := aD.actorRepo.AddActor(&actor)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(resultActor)
}

// swagger:route PUT /actors/{id} Actors updateActor
// Updates actor by id
//
// responses:
//  200: basicResponse
//  400: basicResponse
//	500: basicResponse
func (aD *actorDelivery) UpdateActor(w http.ResponseWriter, r *http.Request) {
	jsonEnc := json.NewEncoder(w)

	id := strings.TrimPrefix(r.URL.Path, "/actors/")
	actorID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	var actor models.Actor
	err = json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	err = aD.actorRepo.UpdateActor(actorID, &actor)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(&models.BasicResponse{Status: "OK"})
}

// swagger:route DELETE /actors/{id} Actors deleteActor
// Deletes actor by id
// responses:
//  200: basicResponse
//  400: basicResponse
//	500: basicResponse
func (aD *actorDelivery) DeleteActor(w http.ResponseWriter, r *http.Request) {
	jsonEnc := json.NewEncoder(w)

	id := strings.TrimPrefix(r.URL.Path, "/actors/")
	actorID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	err = aD.actorRepo.DeleteActor(actorID)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(&models.BasicResponse{Status: "OK"})
}