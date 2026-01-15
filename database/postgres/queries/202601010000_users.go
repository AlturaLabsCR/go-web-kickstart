package queries

import (
	"context"
	"encoding/json"
	"time"

	"app/database/models"
	"app/database/postgres/db"
)

func (pq *PostgresQuerier) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if userStr, err := pq.cache.Get(ctx, userID); err == nil {
		user := &models.User{}
		if err := json.Unmarshal([]byte(userStr), user); err == nil {
			return user, nil
		}
	}

	u, err := pq.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        u.ID,
		CreatedAt: time.Unix(u.CreatedAt, 0),
	}

	if b, err := json.Marshal(user); err == nil {
		_ = pq.cache.Set(ctx, userID, string(b))
	}

	return user, nil
}

func (pq *PostgresQuerier) SetUser(ctx context.Context, id string) error {
	createdAt := time.Now()
	user := &models.User{
		ID:        id,
		CreatedAt: createdAt,
	}

	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := pq.queries.SetUser(ctx, db.SetUserParams{
		ID:        id,
		CreatedAt: createdAt.Unix(),
	}); err != nil {
		return err
	}

	return pq.cache.Set(ctx, id, string(userStr))
}

func (pq *PostgresQuerier) DelUser(ctx context.Context, id string) error {
	if err := pq.queries.DelUser(ctx, id); err != nil {
		return err
	}
	return pq.cache.Del(ctx, id)
}
