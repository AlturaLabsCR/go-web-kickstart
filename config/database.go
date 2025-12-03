package config

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"app/database"
	"app/sessions"
	"app/storage/kv"
	"app/storage/s3"
)

const (
	SqliteDriver   = "sqlite"
	PostgresDriver = "postgresql"

	sqliteMigrations   = "database/sqlite/migrations"
	postgresMigrations = "database/postgres/migrations"
)

type Migrations map[string]embed.FS

func InitDB(migrations Migrations) (database.Database, kv.Store[sessions.Session], kv.Store[s3.Object]) {
	return clients(Config.Database.ConnString, migrations)
}

func clients(connString string, migrations Migrations) (database.Database, kv.Store[sessions.Session], kv.Store[s3.Object]) {
	ctx := context.Background()
	connDriver := getDriver(connString)

	migFS, ok := migrations[connDriver]
	if !ok {
		panic(fmt.Sprintf("no migration folder provided for driver %q", connDriver))
	}

	switch connDriver {
	case SqliteDriver:
		sqlite := initSqlite(connString)
		runMigrations(ctx, sqlite, migFS, sqliteMigrations)
		return sqlite,
			database.NewSqliteSessionStore(sqlite),
			database.NewSqliteObjectStore(sqlite)

	case PostgresDriver:
		pg := initPostgres(ctx, connString)
		runMigrations(ctx, pg, migFS, postgresMigrations)
		return pg,
			database.NewPostgresSessionStore(pg),
			database.NewPostgresObjectStore(pg)
	}

	panic(fmt.Sprintf("invalid db driver: %s", connDriver))
}

func getDriver(connString string) string {
	if strings.Contains(connString, PostgresDriver) {
		return PostgresDriver
	}

	return SqliteDriver
}

func initSqlite(connString string) *database.Sqlite {
	sqlite, err := database.NewSqlite(connString)
	if err != nil {
		abs, errs := filepath.Abs(connString)
		if errs != nil {
			abs = connString
		}
		dir := filepath.Dir(abs)
		panic(fmt.Sprintf(
			"unable to create sqlite connection: %v\ncheck if the folder `%s` exists and has write permissions",
			err, dir,
		))
	}
	return sqlite
}

func initPostgres(ctx context.Context, connString string) *database.Postgres {
	pg, err := database.NewPostgres(ctx, connString)
	if err != nil {
		panic(fmt.Sprintf("unable to create connection pool: %v", err))
	}
	return pg
}

func runMigrations(ctx context.Context, db database.Database, fsys embed.FS, folder string) {
	entries, err := fsys.ReadDir(folder)
	if err != nil {
		panic(fmt.Sprintf("unable to read migrations folder %q: %v", folder, err))
	}

	if len(entries) == 0 {
		panic(fmt.Sprintf("no migration files found in %q", folder))
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}

	if len(files) == 0 {
		panic(fmt.Sprintf("no .sql migration files found in %q", folder))
	}

	sort.Strings(files)

	for _, fname := range files {
		fullpath := path.Join(folder, fname)

		content, err := fsys.ReadFile(fullpath)
		if err != nil {
			panic(fmt.Sprintf("failed to read migration %q: %v", fullpath, err))
		}

		if len(bytes.TrimSpace(content)) == 0 {
			panic(fmt.Sprintf("migration file %q is empty", fullpath))
		}

		if err := db.ExecSQL(ctx, string(content)); err != nil {
			panic(fmt.Sprintf("migration %q failed: %v", fname, err))
		}
	}
}
