package repository

import (
	"context"
	"database/sql"
	"vk-intern_test-case/models"
	"vk-intern_test-case/utils/database"

	actorQueries "vk-intern_test-case/internal/actor/queries"
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
	} else if err != sql.ErrNoRows {
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

func (aR *ActorRepository) UpdateActor(actorID int, actor *models.Actor) (error) {
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


func (aR *ActorRepository) DeleteActor(actorID int) (error) {
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

func (aR *ActorRepository) FindFilmsByActor(actorName string) ([]models.Film, error) {
	transactionCtx := context.Background()
	tx, err := aR.pool.Begin(transactionCtx)
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
	rows, err := tx.Query(transactionCtx, actorQueries.GetFilmsByActor, &actorName)
	if err != nil {
		return []models.Film{}, err
	}

	for rows.Next() {
		film := models.Film{}
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return []models.Film{}, err
		}
		films = append(films, film)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return []models.Film{}, err
	}
	return films, nil
}