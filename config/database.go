package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	pool, err := pgxpool.New(
		context.Background(),
		Environment[EnvConnstr],
	)
	if err != nil {
		panic(fmt.Sprintf("unable to create connection pool: %v", err))
	}

	return pool
}
