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
	ctx := context.Background()

	connString := Config.DB.ConnString

	connDriver := SqliteDriver
	if strings.Contains(connString, PostgresDriver) {
		connDriver = PostgresDriver
	}

	migFS, ok := migrations[connDriver]
	if !ok {
		panic(fmt.Sprintf("no migration folder provided for driver %q", connDriver))
	}

	var (
		conn         database.Database
		sessionStore kv.Store[sessions.Session]
		objectStore  kv.Store[s3.Object]
	)

	switch connDriver {
	case SqliteDriver:
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

		runMigrations(ctx, sqlite, migFS, sqliteMigrations)

		conn = sqlite
		sessionStore = database.NewSqliteSessionStore(sqlite)
		objectStore = database.NewSqliteObjectStore(sqlite)

	case PostgresDriver:
		pg, err := database.NewPostgres(ctx, connString)
		if err != nil {
			panic(fmt.Sprintf("unable to create connection pool: %v", err))
		}

		runMigrations(ctx, pg, migFS, postgresMigrations)

		conn = pg
		sessionStore = database.NewPostgresSessionStore(pg)
		objectStore = database.NewPostgresObjectStore(pg)

	default:
		panic(fmt.Sprintf("invalid db driver: %s", connDriver))
	}

	return conn, sessionStore, objectStore
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
