package repository

import (
	"context"
	"time"
	actorQueries "vk-intern_test-case/internal/actor/queries"
	filmQueries "vk-intern_test-case/internal/film/queries"
	"vk-intern_test-case/models"
	"vk-intern_test-case/utils/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

const logMessage = "film:repository:"

type FilmRepository struct {
	pool database.PgxIface
}

func NewFilmRepository(pool database.PgxIface) *FilmRepository {
	return &FilmRepository{
		pool: pool,
	}
}

func (fR *FilmRepository) AddFilm(filmWithActors *models.FilmWithActors) (*models.Film, error) {
	message := logMessage + "AddFilm:"
	log.Debug(message + "started")
	transactionCtx := context.Background()
	tx, err := fR.pool.Begin(transactionCtx)
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	row := tx.QueryRow(transactionCtx, filmQueries.CreateFilm,
		&filmWithActors.Title,
		&filmWithActors.Description,
		&filmWithActors.ReleaseDate,
		&filmWithActors.Rating,
	)

	err = row.Scan(&filmWithActors.ID)
	if err != nil {
		return nil, err
	}

	for _, actor := range filmWithActors.Actors {
		var actorID int
		row := tx.QueryRow(transactionCtx, actorQueries.GetActorIdByName, &actor)
		err = row.Scan(&actorID)
		if err != nil {
			if err == pgx.ErrNoRows {
				err = nil
				continue
			}
			log.Error(message + err.Error())
			return nil, err
		}

		_, err = tx.Exec(transactionCtx, filmQueries.MakeConnectionFilmWithActor, &actorID, &filmWithActors.Film.ID)
		if err != nil {
			return nil, err
		}
	}

	return &filmWithActors.Film, nil
}

func (fR *FilmRepository) UpdateFilm(filmID int, film *models.Film) error {
	transactionCtx := context.Background()
	tx, err := fR.pool.Begin(transactionCtx)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	_, err = tx.Exec(transactionCtx, filmQueries.UpdateFilm, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating, &filmID)
	if err != nil {
		return err
	}

	return nil
}

func (fR *FilmRepository) DeleteFilm(filmID int) error {
	transactionCtx := context.Background()
	tx, err := fR.pool.Begin(transactionCtx)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	_, err = tx.Exec(transactionCtx, filmQueries.DeleteFilm, &filmID)
	if err != nil {
		return err
	}

	return nil
}

func (fR *FilmRepository) GetFilmsSorted(field string) ([]models.Film, error) {
	transactionCtx := context.Background()
	tx, err := fR.pool.Begin(transactionCtx)
	if err != nil {
		return []models.Film{}, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	query := filmQueries.GetFilmsSortedByRating
	switch field {
	case "release_date":
		query = filmQueries.GetFilmsSortedByDate
	case "title":
		query = filmQueries.GetFilmsSortedByTitle
	}

	films := []models.Film{}
	rows, err := tx.Query(transactionCtx, query)
	if err != nil {
		return []models.Film{}, err
	}

	for rows.Next() {
		film := models.Film{}
		var releaseDatePG pgtype.Date
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &releaseDatePG, &film.Rating)
		if err != nil {
			return []models.Film{}, err
		}
		film.ReleaseDate = releaseDatePG.Time.Format(time.DateOnly)
		films = append(films, film)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return []models.Film{}, err
	}
	return films, nil
}

func (fR *FilmRepository) GetFilmsByTitle(title string) ([]models.Film, error) {
	transactionCtx := context.Background()
	tx, err := fR.pool.Begin(transactionCtx)
	if err != nil {
		return []models.Film{}, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	films := []models.Film{}
	rows, err := tx.Query(transactionCtx, filmQueries.GetFilmsByTitle, &title)
	if err != nil {
		return []models.Film{}, err
	}

	for rows.Next() {
		film := models.Film{}
		var releaseDatePG pgtype.Date
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &releaseDatePG, &film.Rating)
		if err != nil {
			return []models.Film{}, err
		}
		film.ReleaseDate = releaseDatePG.Time.Format(time.DateOnly)
		films = append(films, film)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return []models.Film{}, err
	}
	return films, nil
}

func (fR *FilmRepository) GetFilmsByActor(actorName string) ([]models.Film, error) {
	transactionCtx := context.Background()
	tx, err := fR.pool.Begin(transactionCtx)
	if err != nil {
		return []models.Film{}, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	films := []models.Film{}
	rows, err := tx.Query(transactionCtx, filmQueries.GetFilmsByActor, &actorName)
	if err != nil {
		return []models.Film{}, err
	}

	for rows.Next() {
		film := models.Film{}
		var releaseDatePG pgtype.Date
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &releaseDatePG, &film.Rating)
		if err != nil {
			return []models.Film{}, err
		}
		film.ReleaseDate = releaseDatePG.Time.Format(time.DateOnly)
		films = append(films, film)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return []models.Film{}, err
	}
	return films, nil
}
