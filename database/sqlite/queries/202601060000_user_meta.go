package queries

import (
	"context"

	"app/database/models"
)

func (sq *SqliteQuerier) GetUserMeta(ctx context.Context, userID string) (*models.UserMeta, error) {
	meta, err := sq.queries.GetUserMeta(ctx, userID)
	if err != nil {
		return nil, err
	}

	perms, err := sq.queries.GetPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserMeta{
		ID:      meta.ID,
		Created: meta.Created,
		Name:    meta.Name,
		Perms:   perms,
	}, nil
}
