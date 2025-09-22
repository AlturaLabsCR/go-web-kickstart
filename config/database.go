package config

import (
	"database/sql"

	_ "modernc.org/sqlite"
	// _ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/sijms/go-ora/v2"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbConn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
