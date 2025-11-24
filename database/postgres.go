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

func NewPostgres(ctx context.Context, connString string) (Database, error) {
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

func (db *Postgres) insertUser(ctx context.Context, ownerEmail string) (int64, error) {
	return db.Queries.InsertUser(ctx, ownerEmail)
}

func (db *Postgres) selectUserEmails(ctx context.Context) ([]string, error) {
	return db.Queries.SelectUserEmails(ctx)
}

func (db *Postgres) Close(ctx context.Context) {
	db.Pool.Close()
}
