package postgres

import (
	"app/database"
	"app/database/postgres/queries"
)

var (
	_ database.Database = (*Postgres)(nil)
	_ database.Querier  = (*queries.PostgresQuerier)(nil)
)
