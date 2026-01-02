package queries

import (
	"context"
	"encoding/json"
	"time"

	"app/database/models"
	"app/database/sqlite/db"
)

func (sq *SqliteQuerier) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if userStr, err := sq.cache.Get(ctx, userID); err == nil {
		user := &models.User{}
		if err := json.Unmarshal([]byte(userStr), user); err == nil {
			return user, nil
		}
	}

	u, err := sq.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        u.ID,
		CreatedAt: time.Unix(u.CreatedAt, 0),
	}

	if b, err := json.Marshal(user); err == nil {
		_ = sq.cache.Set(ctx, userID, string(b))
	}

	return user, nil
}

func (sq *SqliteQuerier) SetUser(ctx context.Context, id string, user *models.User) error {
	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := sq.queries.SetUser(ctx, db.SetUserParams{
		ID:        id,
		CreatedAt: user.CreatedAt.Unix(),
	}); err != nil {
		return err
	}

	return sq.cache.Set(ctx, id, string(userStr))
}

func (sq *SqliteQuerier) DelUser(ctx context.Context, id string) error {
	if err := sq.queries.DelUser(ctx, id); err != nil {
		return err
	}
	return sq.cache.Del(ctx, id)
}
