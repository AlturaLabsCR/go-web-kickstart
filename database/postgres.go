package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"app/database/postgres/db"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewPostgres(ctx context.Context, connString string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	queries := db.New(pool)

	return &Postgres{
		Pool:    pool,
		Queries: queries,
	}, nil
}

func (p *Postgres) upsertUser(ctx context.Context, userID string) error {
	return p.Queries.UpsertUser(ctx, userID)
}

func (p *Postgres) Close(ctx context.Context) {
	p.Pool.Close()
}
