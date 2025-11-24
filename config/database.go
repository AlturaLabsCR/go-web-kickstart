package config

import (
	"context"
	"fmt"

	"app/database"
)

func InitDB() database.Database {
	ctx := context.Background()
	connDriver := Environment[EnvDriver]
	connString := Environment[EnvConnStr]

	var conn database.Database
	var err error

	switch connDriver {
	case "sqlite":
		conn, err = database.NewSqlite(connString)
	case "postgres":
		conn, err = database.NewPostgres(ctx, connString)
	default:
		panic(fmt.Sprintf("invalid db driver: %s", connDriver))
	}

	if err != nil {
		panic(fmt.Sprintf("unable to create connection pool: %v", err))
	}

	return conn
}
