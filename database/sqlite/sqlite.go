// Package sqlite
package sqlite

import (
	"context"
	"database/sql"

	"app/cache"
	"app/database"
	"app/database/sqlite/db"
	"app/database/sqlite/queries"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
	db      *sql.DB
	queries *db.Queries
	cache   cache.Cache
}

type SqliteOption func(*Sqlite)

func WithCache(c cache.Cache) SqliteOption {
	return func(s *Sqlite) {
		if c != nil {
			s.cache = c
		}
	}
}

func NewSqlite(connStr string, opts ...SqliteOption) (*Sqlite, error) {
	conn, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	s := &Sqlite{
		db:      conn,
		queries: db.New(conn),
		cache:   database.NoopCache{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

func (s *Sqlite) Exec(ctx context.Context, sql string) error {
	_, err := s.db.ExecContext(ctx, sql)
	return err
}

func (s *Sqlite) Close(_ context.Context) error {
	return s.db.Close()
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

	qtx := queries.New(s.queries.WithTx(tx), s.cache)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit()
}
