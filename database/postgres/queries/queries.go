// Package queries
package queries

import (
	"app/cache"
	"app/database/postgres/db"
)

type PostgresQuerier struct {
	queries *db.Queries
	cache   cache.Cache
}

func New(queries *db.Queries, cache cache.Cache) *PostgresQuerier {
	return &PostgresQuerier{
		queries: queries,
		cache:   cache,
	}
}
