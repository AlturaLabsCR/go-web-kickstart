package queries

import (
	"context"
	"errors"

	"app/database/postgres/db"

	"github.com/jackc/pgx/v5"
)

func (pq *PostgresQuerier) UpsertUserName(ctx context.Context, userName, userID string) error {
	return pq.queries.UpsertUserName(ctx, db.UpsertUserNameParams{
		Name: userName,
		User: userID,
	})
}

func (pq *PostgresQuerier) GetUserName(ctx context.Context, userID string) (string, error) {
	userName, err := pq.queries.GetUserName(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return userName, nil
}
