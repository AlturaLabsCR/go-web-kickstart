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

func (sq *SqliteQuerier) GetUsersMeta(ctx context.Context) ([]models.UserMeta, error) {
	allFromDB, err := sq.queries.GetUsersMeta(ctx)
	if err != nil {
		return nil, err
	}

	var allUsers []models.UserMeta

	for _, u := range allFromDB {
		perms, err := sq.queries.GetPermissions(ctx, u.ID)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, models.UserMeta{
			ID:      u.ID,
			Created: u.Created,
			Name:    u.Name,
			Perms:   perms,
		})
	}

	return allUsers, nil
}
