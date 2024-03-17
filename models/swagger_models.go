package models

// Basic response: Good - OK. Else - Error
// swagger:response basicResponse
type basicResponseWrapper struct {
	// Response Message
	// in: body
	Body BasicResponse
}

// swagger:parameters updateActor
type actorIDParameterWrapper struct {
	// The id of actor
	// in: path
	// required: true
	ID int
}

// An actor from database
// swagger:response actor
type actorResponseWrapper struct {
	// Actor from database
	// in: body
	Body Actor
}

// An actor from database
// swagger:parameters updateActor deleteActor
type actorParameterWrapper struct {
	// Actor from database
	// in: body
	Body Actor
}

// Model for adding film into database
// swagger:parameters addFilm
type filmWithActorsWrapper struct {
	//Film with actors
	// in: body
	Body FilmWithActorsRequest
}

// Model for adding film into database
// swagger:parameters updateFilm
type filmWrapper struct {
	// Film
	// in: body
	Body Film
}

// swagger:parameters updateFilm deleteFilm
type filmIDParameterWrapper struct {
	// The id of film
	// in: path
	// required: true
	ID int `json:"id"`
}
