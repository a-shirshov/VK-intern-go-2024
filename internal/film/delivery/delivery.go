package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"vk-intern_test-case/internal/film"
	"vk-intern_test-case/models"
	"vk-intern_test-case/utils/response"

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

func (fD *FilmDelivery) HandleFilms(w http.ResponseWriter, r *http.Request) {
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
// Добавляет новый фильм в систему, совместно со списком актёров. 
// Если актёра нет в базе - он пропускается и не записывается.
// Актёр добавляется заранее. Поиск происходит по имени.
// security:
// - key:
// responses:
//
//	200: basicResponse
//  401: basicResponse
func (fD *FilmDelivery) AddFilm(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "AddFilm:"
	log.Debug(message + "started")
	jsonEnc := response.MakeJsonEncoder(w)
	var film models.FilmWithActors
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	resultFilm, err := fD.filmRepo.AddFilm(&film)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteResponse(w, jsonEnc, http.StatusOK, resultFilm)
}

// swagger:route PUT /films/{id} Films updateFilm
// Обновляет информацию о фильме, на вход полный поступает вся информация о фильме.
// security:
// - key:
// responses:
//
//	200: basicResponse
//	400: basicResponse
//  401: basicResponse
//	500: basicResponse
func (fD *FilmDelivery) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	jsonEnc := response.MakeJsonEncoder(w)
	var film models.Film
	id := strings.TrimPrefix(r.URL.Path, "/films/")
	filmID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	err = json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}
	err = fD.filmRepo.UpdateFilm(filmID, &film)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteBasicResponse(w, jsonEnc, http.StatusOK, "OK")
}

// swagger:route DELETE /films/{id} Films deleteFilm
// Удаляет фильм из системы
// security:
// - key:
// responses:
//
//	200: basicResponse
//	400: basicResponse
//  401: basicResponse
//	500: basicResponse
func (fD *FilmDelivery) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	jsonEnc := response.MakeJsonEncoder(w)
	id := strings.TrimPrefix(r.URL.Path, "/films/")
	filmID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}
	err = fD.filmRepo.DeleteFilm(filmID)
	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteBasicResponse(w, jsonEnc, http.StatusOK, "OK")
}

// swagger:route GET /films Films getFilms
// Возвращает список фильмов. Можно указать поле для сортировки, по умолчанию - по рейтингу
// responses:
//
//	200: []film
//	500: basicResponse
func (fD *FilmDelivery) GetFilms(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetFilms:"
	log.Debug(message + "started")
	jsonEnc := response.MakeJsonEncoder(w)
	sortBy := r.URL.Query().Get("sort_by")

	resultFilms, err := fD.filmRepo.GetFilmsSorted(sortBy)

	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteResponse(w, jsonEnc, http.StatusOK, resultFilms)
}

// swagger:route GET /film Films getFilm
// Возвращает список фильмов по фрагменту названия или фрагмена имени актёра
// responses:
//
//	200: []film
//	500: basicResponse
func (fD *FilmDelivery) HandleFilm(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		title := r.URL.Query().Get("title")
		if title != "" {
			fD.GetFilmByTitle(w, r)
			return
		}
		actor := r.URL.Query().Get("actor")
		if actor != "" {
			fD.GetFilmByActor(w, r)
			return
		}
		fD.GetFilms(w, r)
	}
}

func (fD *FilmDelivery) GetFilmByTitle(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetFilmByTitle:"
	log.Debug(message + "started")
	jsonEnc := response.MakeJsonEncoder(w)
	title := r.URL.Query().Get("title")

	resultFilms, err := fD.filmRepo.GetFilmsByTitle(title)

	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteResponse(w, jsonEnc, http.StatusOK, resultFilms)
}

func (fD *FilmDelivery) GetFilmByActor(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetFilmByActor:"
	log.Debug(message + "started")
	jsonEnc := response.MakeJsonEncoder(w)
	actor := r.URL.Query().Get("actor")

	resultFilms, err := fD.filmRepo.GetFilmsByActor(actor)

	if err != nil {
		log.Error(err)
		response.WriteBasicResponse(w, jsonEnc, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteResponse(w, jsonEnc, http.StatusOK, resultFilms)
}
