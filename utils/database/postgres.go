package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), "postgresql://artyom:123456@postgres:5432/filmbase")
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return pool, nil
}
