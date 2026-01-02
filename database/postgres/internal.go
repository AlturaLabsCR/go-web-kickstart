package postgres

import "app/database"

var (
	_ database.Database = (*Postgres)(nil)
	_ database.Querier  = (*PostgresQuerier)(nil)
)
