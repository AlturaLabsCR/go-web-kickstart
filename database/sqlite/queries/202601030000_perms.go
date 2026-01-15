package queries

import (
	"context"

	"app/database/sqlite/db"
)

func (sq *SqliteQuerier) GetRoles(ctx context.Context, userID string) ([]string, error) {
	return sq.queries.GetRoles(ctx, userID)
}

func (sq *SqliteQuerier) GetPermissions(ctx context.Context, userID string) ([]string, error) {
	return sq.queries.GetPermissions(ctx, userID)
}

func (sq *SqliteQuerier) SetRole(ctx context.Context, userID string, roleName string) error {
	return sq.queries.SetRole(ctx, db.SetRoleParams{
		User: userID,
		Role: roleName,
	})
}
