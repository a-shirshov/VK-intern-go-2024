package queries

const (
	CreateFilm = `insert into film (name, description, release_date, rating)
		values ($1, $2, $3, $4) returning id;`
	MakeConnectionFilmWithActor = `insert into actor_film (actor_id, film_id) values ($1, $2);`
	UpdateFilm                  = `update film set name = $1, description = $2, release_date = $3, rating = $4 where id = $5;`
	DeleteFilm                  = `delete from film where id = $1`
	GetFilmsSortedByRating      = `select * from film order by rating desc;`
	GetFilmsSortedByDate        = `select * from film order by release_date;`
	GetFilmsSortedByTitle       = `select * from film order by title;`
	GetFilmsByTitle             = `select * from film where title = $1%;`
)
