package repository

import (
	"testing"
	"vk-intern_test-case/models"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func prepareTestEnvironment(t *testing.T) (*ActorRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	testActorRepo := NewActorRepository(mock)
	return testActorRepo, mock
}

func TestShouldSuccessfullyAddNewActor(t *testing.T) {
	actorRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	newActor := &models.Actor{
		ActorRequest: models.ActorRequest{
			Name:        "Леонардо Ди Каприо",
			Gender:      "Мужской",
			DateOfBirth: "2001-08-06",
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery("select id from actor").WithArgs(&newActor.Name).WillReturnError(pgx.ErrNoRows)
	mock.ExpectQuery("insert into actor").WithArgs(&newActor.Name, &newActor.Gender, &newActor.DateOfBirth).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	resultActor, err := actorRepo.AddActor(newActor)
	if err != nil {
		t.Errorf("error was not expected while adding an actor: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, 1, resultActor.ID)
	assert.Nil(t, err)
}

func TestShouldReturnExistingActorIfAddingTheSame(t *testing.T) {
	actorRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	newActor := &models.Actor{
		ActorRequest: models.ActorRequest{
			Name:        "Actor_name",
			Gender:      "Мужской",
			DateOfBirth: "06-08-2001",
		},
	}
	mock.ExpectBegin()
	mock.ExpectQuery("select id from actor").WithArgs(&newActor.Name).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(5))
	mock.ExpectCommit()

	resultActor, err := actorRepo.AddActor(newActor)
	if err != nil {
		t.Errorf("error was not expected while adding an actor: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, 5, resultActor.ID)
	assert.Nil(t, err)
}

func TestShouldSuccessfullyUpdateAnExistingActor(t *testing.T) {
	actorRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	newActor := &models.Actor{
		ActorRequest: models.ActorRequest{
			Name:        "Actor_name",
			Gender:      "Мужской",
			DateOfBirth: "06-08-2001",
		},
	}
	actorID := 5

	mock.ExpectBegin()
	mock.ExpectExec("update actor").WithArgs(&newActor.Name, &newActor.Gender, &newActor.DateOfBirth, &actorID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	err := actorRepo.UpdateActor(actorID, newActor)
	if err != nil {
		t.Errorf("error was not expected while adding an actor: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Nil(t, err)
}

func TestShouldSuccessfullyDeleteAnExistingActor(t *testing.T) {
	actorRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	actorID := 5

	mock.ExpectBegin()
	mock.ExpectExec("delete from actor").WithArgs(&actorID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	err := actorRepo.DeleteActor(actorID)
	if err != nil {
		t.Errorf("error was not expected while adding an actor: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Nil(t, err)
}

func TestShouldSuccessfullyGetActors(t *testing.T) {
	actorRepo, mock := prepareTestEnvironment(t)
	defer mock.Close()
	actorIDs := []int{1, 2}
	mock.ExpectBegin()
	mock.ExpectQuery("select").
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "gender", "date_of_birth"}).
		AddRow(actorIDs[0], "Леонардо Ди Каприо", "Мужской", "2024-03-18").
		AddRow(actorIDs[1], "Марго Робби", "Женский", "2023-03-18")).
		RowsWillBeClosed()
	for i := 0; i < len(actorIDs); i++ {
		mock.ExpectQuery("select").WithArgs(&actorIDs[i]).WillReturnRows().
			WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "release_date", "rating"}).
			AddRow(1, "Титаник", "cool film", "2020-06-10", 8).
			AddRow(2, "Не титаник", "super film", "2018-06-10", 7)).
			RowsWillBeClosed()
	}
	mock.ExpectCommit()

	_, err := actorRepo.GetActors()
	if err != nil {
		t.Errorf("error was not expected while adding an actor: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Nil(t, err)
}