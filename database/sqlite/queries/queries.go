// Package queries
package queries

import (
	"app/cache"
	"app/database/sqlite/db"
)

type SqliteQuerier struct {
	queries *db.Queries
	cache   cache.Cache
}

func New(queries *db.Queries, cache cache.Cache) *SqliteQuerier {
	return &SqliteQuerier{
		queries: queries,
		cache:   cache,
	}
}
