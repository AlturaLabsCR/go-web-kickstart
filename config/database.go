package config

import (
	"context"
	"fmt"

	"app/database"
	"app/sessions"
	"app/storage/kv"
)

func InitDB() (database.Database, kv.Store[sessions.Session]) {
	ctx := context.Background()
	connDriver := Environment[EnvDriver]
	connString := Environment[EnvConnStr]

	var conn database.Database
	var store kv.Store[sessions.Session]

	switch connDriver {
	case "sqlite":
		sqlite, err := database.NewSqlite(connString)
		if err != nil {
			panic(fmt.Sprintf("unable to create sqlite connection: %v", err))
		}
		sqlitekv := database.NewSqliteSessionStore(sqlite)

		conn = sqlite
		store = sqlitekv
	case "postgres":
		pg, err := database.NewPostgres(ctx, connString)
		if err != nil {
			panic(fmt.Sprintf("unable to create connection pool: %v", err))
		}
		pgkv := database.NewPostgresSessionStore(pg)

		conn = pg
		store = pgkv
	default:
		panic(fmt.Sprintf("invalid db driver: %s", connDriver))
	}

	return conn, store
}
