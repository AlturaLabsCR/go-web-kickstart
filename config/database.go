package config

import (
	"database/sql"

	_ "modernc.org/sqlite"
	// _ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/sijms/go-ora/v2"
)

const (
	defaultDriver = "sqlite"
	defaultConn   = "./db.db"
)

func InitDB(d, c string) (*sql.DB, error) {
	if d == "" {
		d = defaultDriver
	}

	if c == "" {
		c = defaultConn
	}

	db, err := sql.Open(d, c)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
