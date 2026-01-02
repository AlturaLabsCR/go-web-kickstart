package sqlite

import (
	"app/database"
	"app/database/sqlite/queries"
)

var (
	_ database.Database = (*Sqlite)(nil)
	_ database.Querier  = (*queries.SqliteQuerier)(nil)
)
