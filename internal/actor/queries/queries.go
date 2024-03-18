package queries

const (
	CreateAnActor     = `insert into actor (name, gender, date_of_birth) values ($1, $2, $3) returning id;`
	GetActorIdByName  = `select id from actor where name = $1;`
	UpdateActor       = `update actor set name = $1, gender = $2, date_of_birth = $3 where id = $4;`
	GetActorByID      = `select * from actor where id = $1;`
	DeleteActor       = `delete from actor where id = $1;`
	GetActors         = `select * from actor;`
	GetFilmsByActorId = `select f.id, f.title, f.description, f.release_date, f.rating from film as f
		join actor_film as af on f.id = af.film_id
		join actor as a on a.id = af.actor_id
		where a.id = $1;`
)
