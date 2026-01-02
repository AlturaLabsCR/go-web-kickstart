package sqlite

import (
	"context"

	"app/cache"
	"app/database"
	"app/database/sqlite/db"
)

type SqliteQuerier struct {
	queries *db.Queries
	cache   cache.Cache
}

func (s *Sqlite) WithTx(
	_ context.Context,
	fn func(q database.Querier) error,
) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := &SqliteQuerier{
		queries: s.queries.WithTx(tx),
		cache:   s.cache,
	}

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit()
}
