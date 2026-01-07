package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"app/database"
	"app/database/postgres"
	"app/database/sqlite"
)

const (
	SqliteDriver   = "sqlite"
	PostgresDriver = "postgresql"
)

func InitDB() (database.Database, error) {
	var empty database.Database

	ctx := context.Background()
	connStr := Config.Database.ConnString

	switch getDriver(connStr) {
	case SqliteDriver:
		if info, err := os.Stat(filepath.Dir(connStr)); err == nil {
			if !info.IsDir() {
				return empty, fmt.Errorf("%s, '%s' %s", connStr, info.Name(), "is not a directory")
			}
		} else {
			return empty, fmt.Errorf("%s, %w", connStr, err)
		}

		sq, err := sqlite.NewSqlite(connStr)
		if err != nil {
			return empty, err
		}
		return sq, nil
	case PostgresDriver:
		pq, err := postgres.NewPostgres(ctx, connStr)
		if err != nil {
			return empty, err
		}
		return pq, nil
	}

	return empty, fmt.Errorf("no available driver for: %s", connStr)
}

func getDriver(connStr string) string {
	if strings.Contains(connStr, PostgresDriver) {
		return PostgresDriver
	}

	return SqliteDriver
}
