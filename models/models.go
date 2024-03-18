package models

type BasicResponse struct {
	Status string `json:"status"`
}

// Actor represents actor in system
// swagger:model actor
type Actor struct {
	// The id for this actor
	//
	// required: true
	// min: 1
	ID int `json:"id"`
	ActorRequest
}

type ActorRequest struct {
	// Name of the actor
	//
	// required: true
	// example: Леонардо Ди Каприо
	Name string `json:"name"`
	// Gender of the actor
	//
	// required: true
	// example: Мужской
	Gender string `json:"gender"`
	// Date of birth of the actor
	//
	// required: true
	// В формате YYYY-MM_DD
	// example: 2001-08-06
	DateOfBirth string `json:"date_of_birth"`
}

type FilmRequest struct {
	// Name of the actor
	//
	// required: true
	// example: Titanic
	Title string `json:"title"`
	// Description of film
	//
	// required: true
	// example: film_description
	Description string `json:"description"`
	// Release date of film
	//
	// required: true
	// В формате YYYY-MM_DD
	// example: 2023-03-17
	ReleaseDate string `json:"release_date"`
	// Rating of the film
	//
	// required: true
	// minimun: 0
	// maximum: 10
	// example: 7
	Rating int `json:"rating"`
}

// Film represents film in system
// swagger:model film
type Film struct {
	// The id for this film
	//
	// min: 1
	ID int `json:"id"`
	FilmRequest
}

type FilmWithActors struct {
	Film
	Actors []string `json:"actors"`
}

// Actor with films in which playing
// swagger:model actorWithFilms
type ActorWithFilms struct {
	Actor
	Films []Film `json:"films"`
}
