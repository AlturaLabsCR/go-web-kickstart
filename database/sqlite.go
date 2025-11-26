package database

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	"app/database/sqlite/db"
)

type Sqlite struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewSqlite(connString string) (*Sqlite, error) {
	conn, err := sql.Open("sqlite", connString)
	if err != nil {
		return nil, err
	}

	queries := db.New(conn)

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Sqlite{
		DB:      conn,
		Queries: queries,
	}, nil
}

func (s *Sqlite) upsertUser(ctx context.Context, userID string) error {
	return s.Queries.UpsertUser(ctx, userID)
}

func (s *Sqlite) Close(ctx context.Context) {
	s.DB.Close()
}
