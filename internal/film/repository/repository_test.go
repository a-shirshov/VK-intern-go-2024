package repository

import (
	"testing"
	"vk-intern_test-case/models"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func prepareTestEnvironment(t *testing.T) (*FilmRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	testFilmRepo := NewFilmRepository(mock)
	return testFilmRepo, mock
}

func TestShouldSuccessfullyAddNewFilm(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	newFilm := &models.FilmWithActors{
		Film: models.Film{
			FilmRequest: models.FilmRequest{
				Title:       "film_title",
				Description: "film_description",
				ReleaseDate: "06-08-2001",
				Rating:      8,
			},
		},
		Actors: []string{"Leo Di", "Keanu Rea"},
	}
	newFilmID := 1

	mock.ExpectBegin()
	mock.ExpectQuery("insert into film").WithArgs(&newFilm.Title, &newFilm.Description, &newFilm.ReleaseDate, &newFilm.Rating).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newFilmID))
	for index := range newFilm.Actors {
		actorID := index + 1
		mock.ExpectQuery("select id from actor").WithArgs(&newFilm.Actors[index]).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(actorID))
		mock.ExpectExec("insert into actor_film").WithArgs(&actorID, &newFilmID).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}
	mock.ExpectCommit()

	resultFilm, err := filmRepo.AddFilm(newFilm)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, 1, resultFilm.ID)
	assert.Equal(t, newFilm.FilmRequest, resultFilm.FilmRequest)
	assert.Nil(t, err)
}

func TestShouldSuccessfullyAddNewFilmWithUnknownActors(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	newFilm := &models.FilmWithActors{
		Film: models.Film{
			FilmRequest: models.FilmRequest{
				Title:       "film_title",
				Description: "film_description",
				ReleaseDate: "06-08-2001",
				Rating:      8,
			},
		},
		Actors: []string{"Leo Di", "Keanu Rea"},
	}
	newFilmID := 1

	mock.ExpectBegin()
	mock.ExpectQuery("insert into film").WithArgs(&newFilm.Title, &newFilm.Description, &newFilm.ReleaseDate, &newFilm.Rating).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newFilmID))
	for index := range newFilm.Actors {
		mock.ExpectQuery("select id from actor").WithArgs(&newFilm.Actors[index]).WillReturnError(pgx.ErrNoRows)
	}
	mock.ExpectCommit()

	resultFilm, err := filmRepo.AddFilm(newFilm)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, 1, resultFilm.ID)
	assert.Equal(t, newFilm.FilmRequest, resultFilm.FilmRequest)
	assert.Nil(t, err)
}

func TestShouldSuccessfullyUpdateFilm(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	newFilm := &models.Film{
		FilmRequest: models.FilmRequest{
			Title:       "film_title",
			Description: "film_description",
			ReleaseDate: "06-08-2001",
			Rating:      8,
		},
	}

	newFilmID := 1

	mock.ExpectBegin()
	mock.ExpectExec("update film").WithArgs(&newFilm.Title, &newFilm.Description, &newFilm.ReleaseDate, &newFilm.Rating, &newFilmID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	err := filmRepo.UpdateFilm(newFilmID, newFilm)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Nil(t, err)
}

func TestShouldSuccessfullyDeleteFilm(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	filmID := 1

	mock.ExpectBegin()
	mock.ExpectExec("delete from film").WithArgs(&filmID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	err := filmRepo.DeleteFilm(filmID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Nil(t, err)
}

func TestShouldSuccessfullyReturnFilmsSortedByRating(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	field := "rating"

	mock.ExpectBegin()
	mock.ExpectQuery("select").
		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "release_date", "rating"}).
			AddRow(1, "Titanic", "cool", "2001-08-06", 8).
			AddRow(2, "Titanic 2", "not cool", "2001-08-06", 7)).
		RowsWillBeClosed()
	mock.ExpectCommit()

	resultFilms, err := filmRepo.GetFilmsSorted(field)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, 2, len(resultFilms))
	assert.Nil(t, err)
}

func TestShouldSuccessfullyReturnFilmsByTitle(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	title := "Tit"

	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(&title).
		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "release_date", "rating"}).
			AddRow(1, "Titanic", "cool", "2001-08-06", 8).
			AddRow(2, "Titanic 2", "not cool", "2001-08-06", 7)).
		RowsWillBeClosed()
	mock.ExpectCommit()

	resultFilms, err := filmRepo.GetFilmsByTitle(title)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, 2, len(resultFilms))
	assert.Nil(t, err)
}

func TestShouldSuccessfullyReturnFilmsByActor(t *testing.T) {
	filmRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	actor := "Leo"

	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(&actor).
		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "release_date", "rating"}).
			AddRow(1, "Titanic", "cool", "2001-08-06", 8).
			AddRow(2, "Titanic 2", "not cool", "2001-08-06", 7)).
		RowsWillBeClosed()
	mock.ExpectCommit()

	resultFilms, err := filmRepo.GetFilmsByActor(actor)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, 2, len(resultFilms))
	assert.Nil(t, err)
}
