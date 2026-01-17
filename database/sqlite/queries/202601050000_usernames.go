package queries

import (
	"context"
	"database/sql"
	"errors"

	"app/database/sqlite/db"
)

func (sq *SqliteQuerier) UpsertUserName(ctx context.Context, userName, userID string) error {
	return sq.queries.UpsertUserName(ctx, db.UpsertUserNameParams{
		Name: userName,
		User: userID,
	})
}

func (sq *SqliteQuerier) GetUserName(ctx context.Context, userID string) (string, error) {
	userName, err := sq.queries.GetUserName(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return userName, nil
}
