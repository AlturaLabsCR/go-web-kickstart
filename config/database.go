package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	ctx := context.Background()
	connString := Environment[EnvConnstr]

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(fmt.Sprintf("unable to create connection pool: %v", err))
	}

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("unable to ping connection pool: %v", err))
	}

	return pool
}
