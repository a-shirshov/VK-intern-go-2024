package repository

import (
	"context"
	"time"
	"vk-intern_test-case/models"
	"vk-intern_test-case/utils/database"

	actorQueries "vk-intern_test-case/internal/actor/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

const logMessage = "actor:repository:"

type ActorRepository struct {
	pool database.PgxIface
}

func NewActorRepository(pool database.PgxIface) *ActorRepository {
	return &ActorRepository{
		pool: pool,
	}
}

func (aR *ActorRepository) AddActor(actor *models.Actor) (*models.Actor, error) {
	message := logMessage + "AddActor:"
	log.Debug(message + "started")
	transactionCtx := context.Background()
	tx, err := aR.pool.Begin(transactionCtx)
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

	row := tx.QueryRow(transactionCtx, actorQueries.GetActorIdByName,
		&actor.Name,
	)

	err = row.Scan(&actor.ID)
	if err == nil {
		return actor, nil
	} else if err != pgx.ErrNoRows {
		log.Error(message + err.Error())
		return nil, err
	}

	row = tx.QueryRow(transactionCtx, actorQueries.CreateAnActor,
		&actor.Name,
		&actor.Gender,
		&actor.DateOfBirth,
	)
	err = row.Scan(&actor.ID)
	if err != nil {
		return nil, err
	}
	return actor, nil
}

func (aR *ActorRepository) UpdateActor(actorID int, actor *models.Actor) error {
	transactionCtx := context.Background()
	tx, err := aR.pool.Begin(transactionCtx)
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

	_, err = tx.Exec(transactionCtx, actorQueries.UpdateActor, &actor.Name, &actor.Gender, &actor.DateOfBirth, &actorID)
	if err != nil {
		return err
	}

	return nil
}

func (aR *ActorRepository) DeleteActor(actorID int) error {
	transactionCtx := context.Background()
	tx, err := aR.pool.Begin(transactionCtx)
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

	_, err = tx.Exec(transactionCtx, actorQueries.DeleteActor, &actorID)
	if err != nil {
		return err
	}

	return nil
}

func (aR *ActorRepository) GetActors() ([]models.ActorWithFilms, error) {
	transactionCtx := context.Background()
	tx, err := aR.pool.Begin(transactionCtx)
	if err != nil {
		return []models.ActorWithFilms{}, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(transactionCtx)
		default:
			_ = tx.Rollback(transactionCtx)
		}
	}()

	actorsWithFilms := []models.ActorWithFilms{}
	rows, err := tx.Query(transactionCtx, actorQueries.GetActors)
	if err != nil {
		return []models.ActorWithFilms{}, err
	}

	var actors []models.Actor
	for rows.Next() {
		var actor models.Actor
		var dateOfBirthPG pgtype.Date
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &dateOfBirthPG)
		if err != nil {
			return []models.ActorWithFilms{}, err
		}
		actor.DateOfBirth = dateOfBirthPG.Time.Format(time.DateOnly)
		actors = append(actors, actor)
	}
	rows.Close()

	for _, actor := range actors {
		filmsRows, err := tx.Query(transactionCtx, actorQueries.GetFilmsByActorId, &actor.ID)
		if err != nil {
			return []models.ActorWithFilms{}, err
		}
		var films []models.Film
		for filmsRows.Next() {
			var film models.Film
			var releaseDatePG pgtype.Date
			err := filmsRows.Scan(&film.ID, &film.Title, &film.Description, &releaseDatePG, &film.Rating)
			if err != nil {
				return []models.ActorWithFilms{}, err
			}
			film.ReleaseDate = releaseDatePG.Time.Format(time.DateOnly)
			films = append(films, film)
		}
		actorsWithFilms = append(actorsWithFilms, models.ActorWithFilms{Actor: actor, Films: films})
		filmsRows.Close()
	}

	return actorsWithFilms, nil
}
