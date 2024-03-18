package queries

const (
	CreateFilm = `insert into film (title, description, release_date, rating)
		values ($1, $2, $3, $4) returning id;`
	MakeConnectionFilmWithActor = `insert into actor_film (actor_id, film_id) values ($1, $2);`
	UpdateFilm                  = `update film set title = $1, description = $2, release_date = $3, rating = $4 where id = $5;`
	DeleteFilm                  = `delete from film where id = $1`
	GetFilmsSortedByRating      = `select * from film order by rating desc;`
	GetFilmsSortedByDate        = `select * from film order by release_date;`
	GetFilmsSortedByTitle       = `select * from film order by title;`
	GetFilmsByTitle             = `select * from film where lower(title) like lower($1) || '%';`
	GetFilmsByActor             = `select f.id, f.title, f.description, f.release_date, f.rating from film as f
		join actor_film as af on f.id = af.id
		join actor as a on a.id = af.id
		where lower(a.name) like lower($1) || '%';`
)
