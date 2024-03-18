package models

// Ответ системы. В случае успеха - ОК. Иначе описание ошибки
// swagger:response basicResponse
type basicResponseWrapper struct {
	// Response Message
	// in: body
	Body BasicResponse
}

// swagger:parameters updateActor deleteActor
type actorIDParameterWrapper struct {
	// ID актёра
	// in: path
	// required: true
	ID int `json:"id"`
}

// An actor from database
// swagger:response actor
type actorResponseWrapper struct {
	// Данные об актёре
	// in: body
	Body Actor
}

// An actor from database
// swagger:parameters updateActor
type actorParameterWrapper struct {
	// Данные об актёре
	// in: body
	Body Actor
}

// To add film in a system. Film and actors list
// swagger:model filmWithActors
type FilmWithActorsRequest struct {
	FilmRequest
	Actors []string `json:"actors"`
}

// Model for adding film into database
// swagger:parameters addFilm
type filmWithActorsWrapper struct {
	// Данные об фильме с актёрами
	// in: body
	Body FilmWithActorsRequest
}

// Model for adding actor into database
// swagger:parameters addActor
type actorRequestWrapper struct {
	// Данные о актёре
	// in: body
	Body ActorRequest
}

// Model for adding film into database
// swagger:parameters updateFilm
type filmWrapper struct {
	// Данные о фильме
	// in: body
	Body Film
}

// swagger:parameters updateFilm deleteFilm
type filmIDParameterWrapper struct {
	// ID фильма
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:parameters getFilms
type filmSortByParameterWrapper struct {
	// Параметр для сортировки. Возможные поля - rating, title, release_date
	// in: query
	SortBy string `json:"sort_by"`
}

// swagger:parameters getFilm
type filmByTitleParameterWrapper struct {
	// Поиск по фрагменту названия
	// in: query
	Title string `json:"title"`
}

// swagger:parameters getFilm
type filmByActorNameParameterWrapper struct {
	// Поиск по фрагменту имени актёра
	// in: query
	Actor string `json:"actor"`
}