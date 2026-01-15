package queries

import (
	"context"

	"app/database/postgres/db"
)

func (pq *PostgresQuerier) GetRoles(ctx context.Context, userID string) ([]string, error) {
	return pq.queries.GetRoles(ctx, userID)
}

func (pq *PostgresQuerier) GetPermissions(ctx context.Context, userID string) ([]string, error) {
	return pq.queries.GetPermissions(ctx, userID)
}

func (pq *PostgresQuerier) SetRole(ctx context.Context, userID string, roleName string) error {
	return pq.queries.SetRole(ctx, db.SetRoleParams{
		User: userID,
		Role: roleName,
	})
}
