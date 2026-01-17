package queries

import (
	"context"

	"app/database/models"
)

func (pq *PostgresQuerier) GetUserMeta(ctx context.Context, userID string) (*models.UserMeta, error) {
	meta, err := pq.queries.GetUserMeta(ctx, userID)
	if err != nil {
		return nil, err
	}

	perms, err := pq.queries.GetPermissions(ctx, userID)
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

func (pq *PostgresQuerier) GetUsersMeta(ctx context.Context) ([]models.UserMeta, error) {
	allFromDB, err := pq.queries.GetUsersMeta(ctx)
	if err != nil {
		return nil, err
	}

	var allUsers []models.UserMeta

	for _, u := range allFromDB {
		perms, err := pq.queries.GetPermissions(ctx, u.ID)
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
