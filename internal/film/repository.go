package film

import "vk-intern_test-case/models"

type FilmRepository interface {
	AddFilm(*models.FilmWithActors) (*models.Film, error)
	UpdateFilm(int, *models.Film) (error)
	DeleteFilm(int) (error)
	GetFilmsSortedByTitle() ([]models.Film, error)
	GetFilmsSortedByDate() ([]models.Film, error)
	GetFilmsSortedByRating() ([]models.Film, error)
}