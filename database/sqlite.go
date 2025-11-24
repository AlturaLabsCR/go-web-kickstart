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

func NewSqlite(connString string) (Database, error) {
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

func (db *Sqlite) InsertOwner(ctx context.Context, ownerName string) (int64, error) {
	return db.Queries.InsertOwner(ctx, ownerName)
}
