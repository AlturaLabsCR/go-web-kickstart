// Package postgres
package postgres

import (
	"context"

	"app/cache"
	"app/database"
	"app/database/postgres/db"
	"app/database/postgres/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db      *pgxpool.Pool
	queries *db.Queries
	cache   cache.Cache
}

type PostgresOption func(*Postgres)

func WithCache(c cache.Cache) PostgresOption {
	return func(p *Postgres) {
		if c != nil {
			p.cache = c
		}
	}
}

func NewPostgres(
	ctx context.Context,
	connStr string,
	opts ...PostgresOption,
) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	p := &Postgres{
		db:      pool,
		queries: db.New(pool),
		cache:   database.NoopCache{},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

func (p *Postgres) Exec(ctx context.Context, sql string) error {
	_, err := p.db.Exec(ctx, sql)
	return err
}

func (p *Postgres) Close(ctx context.Context) error {
	p.db.Close()
	return nil
}

func (p *Postgres) WithTx(
	ctx context.Context,
	fn func(q database.Querier) error,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := queries.New(p.queries.WithTx(tx), p.cache)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
