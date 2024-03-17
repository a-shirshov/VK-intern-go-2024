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
	Name string `json:"name"`
	// Gender of the actor
	//
	// required: true
	Gender string `json:"gender"`
	// Date of birth of the actor
	//
	// required: true
	DateOfBirth string `json:"dateOfBirth"`
}

type FilmRequest struct {
	// Name of the actor
	//
	// required: true
	Title string `json:"title"`
	// Description of film
	//
	// required: true
	Description string `json:"description"`
	// Release date of film
	//
	// required: true
	// example: 2023-03-17
	ReleaseDate string `json:"release_date"`
	// Rating of the film
	//
	// required: true
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

// To add film in a system. Film and actors list
// swagger:model filmWithActors
type FilmWithActorsRequest struct {
	FilmRequest
	Actors []string `json:"actors"`
}

type FilmWithActors struct {
	Film
	Actors []string `json:"actors"`
}