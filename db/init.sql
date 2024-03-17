CREATE TABLE IF NOT EXISTS actor (
    id serial not null unique,
    name text not null unique,
    gender text not null,
    date_of_birth date not null 
);

CREATE TABLE IF NOT EXISTS film (
    id serial not null unique,
    name text not null,
    description text not null,
    release_date date not null,
    rating smallint not null
);

CREATE TABLE IF NOT EXISTS actor_film (
    id serial not null,
    actor_id int,
    film_id int,
    FOREIGN KEY (actor_id) REFERENCES actor (id) on delete cascade,
    FOREIGN KEY (film_id) REFERENCES film (id) on delete cascade
);