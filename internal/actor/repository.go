package actor

import "vk-intern_test-case/models"

type ActorRepository interface {
	AddActor(*models.Actor) (*models.Actor, error)
	UpdateActor(int, *models.Actor) (error)
	DeleteActor(int) (error)
}