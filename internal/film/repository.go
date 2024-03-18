package film

import "vk-intern_test-case/models"

type FilmRepository interface {
	AddFilm(*models.FilmWithActors) (*models.Film, error)
	UpdateFilm(int, *models.Film) error
	DeleteFilm(int) error
	GetFilmsSorted(field string) ([]models.Film, error)
	GetFilmsByTitle(title string) ([]models.Film, error)
	GetFilmsByActor(actorName string) ([]models.Film, error)
}
