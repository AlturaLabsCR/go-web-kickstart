package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"app/cache"
	"app/database"
	"app/database/postgres"
	"app/database/sqlite"
)

const (
	SqliteDriver   = "sqlite"
	PostgresDriver = "postgresql"
)

func InitDB(ctx context.Context) (database.Database, error) {
	var empty database.Database

	connStr := Config.Database.ConnString

	driver, err := getDriver(connStr)
	if err != nil {
		return empty, fmt.Errorf("no available driver: %w", err)
	}

	switch driver {
	case SqliteDriver:
		return initSqlite(connStr)
	case PostgresDriver:
		return initPostgres(ctx, connStr)
	}

	return empty, fmt.Errorf("no available driver for: %s", connStr)
}

func initSqlite(connStr string) (*sqlite.Sqlite, error) {
	return sqlite.NewSqlite(connStr, sqlite.WithCache(cache.NewMemoryStore()))
}

func initPostgres(ctx context.Context, connStr string) (*postgres.Postgres, error) {
	return postgres.NewPostgres(ctx, connStr, postgres.WithCache(cache.NewMemoryStore()))
}

func getDriver(connStr string) (string, error) {
	if strings.Contains(connStr, PostgresDriver) {
		return PostgresDriver, nil
	}

	if info, err := os.Stat(filepath.Dir(connStr)); err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("%s, '%s' %s", connStr, info.Name(), "is not a directory")
		}
	} else {
		return "", fmt.Errorf("%s, %w", connStr, err)
	}

	return SqliteDriver, nil
}
