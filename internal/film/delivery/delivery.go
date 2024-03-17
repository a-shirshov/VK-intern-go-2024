package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"vk-intern_test-case/internal/film"
	"vk-intern_test-case/models"

	log "github.com/sirupsen/logrus"
)

const logMessage = "film:delivery:"

type FilmDelivery struct {
	filmRepo film.FilmRepository
}

func NewFilmDelivery(fR film.FilmRepository) *FilmDelivery {
	return &FilmDelivery{
		filmRepo: fR,
	}
}

func (fD *FilmDelivery) HandleFilm(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "HandleFilm:"
	log.Debug(message + "started")
	switch r.Method {
	case http.MethodGet:
		fD.GetFilms(w, r)
	case http.MethodPost:
		fD.AddFilm(w, r)
	case http.MethodPut:
		fD.UpdateFilm(w, r)
	case http.MethodDelete:
		fD.DeleteFilm(w, r)
	}
}

// swagger:route POST /films Films addFilm
// Returns a status of operation adding an film
// responses:
//
//  200: basicResponse
func (fD *FilmDelivery) AddFilm(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "AddFilm:"
	log.Debug(message + "started")
	var film models.FilmWithActors
	jsonEnc := json.NewEncoder(w)
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	resultFilm, err := fD.filmRepo.AddFilm(&film)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(resultFilm)
}

// swagger:route PUT /films/{id} Films updateFilm
// Updates film by id
// responses:
//
//  200: basicResponse
//  400: basicResponse
//  500: basicResponse
func (fD *FilmDelivery) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	jsonEnc := json.NewEncoder(w)
	var film models.Film
	id := strings.TrimPrefix(r.URL.Path, "/films/")
	filmID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	err = fD.filmRepo.UpdateFilm(filmID, &film)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(&models.BasicResponse{Status: "OK"})
}

// swagger:route DELETE /films/{id} Films deleteFilm
// Updates film by id
// responses:
//
//  200: basicResponse
//	400: basicResponse
//	500: basicResponse
func (fD *FilmDelivery) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/films/")
	jsonEnc := json.NewEncoder(w)
	filmID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}
	err = fD.filmRepo.DeleteFilm(filmID)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(&models.BasicResponse{Status: "OK"})
}

func (fD *FilmDelivery) GetFilms(w http.ResponseWriter, r *http.Request) {
	jsonEnc := json.NewEncoder(w)
	sortBy := r.URL.Query().Get("sortBy")
	var resultFilms []models.Film
	var err error
	switch sortBy {
	case "title":
		resultFilms, err = fD.filmRepo.GetFilmsSortedByTitle()
	case "release_date":
		resultFilms, err = fD.filmRepo.GetFilmsSortedByDate()
	default:
		resultFilms, err = fD.filmRepo.GetFilmsSortedByRating()
	}

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(&models.BasicResponse{Status: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(resultFilms)
}
