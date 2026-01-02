package sqlite

import "app/database"

var (
	_ database.Database = (*Sqlite)(nil)
	_ database.Querier  = (*SqliteQuerier)(nil)
)
