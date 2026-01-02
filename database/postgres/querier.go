package postgres

import (
	"context"

	"app/cache"
	"app/database"
	"app/database/postgres/db"
)

type PostgresQuerier struct {
	queries *db.Queries
	cache   cache.Cache
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

	qtx := &PostgresQuerier{
		queries: p.queries.WithTx(tx),
		cache:   p.cache,
	}

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
