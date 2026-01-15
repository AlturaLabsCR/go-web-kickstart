// Package database
package database

import (
	"context"

	"app/database/models"
)

type Database interface {
	Querier() Querier
	WithTx(ctx context.Context, fn func(q Querier) error) (err error)
	Exec(ctx context.Context, sql string) (err error)
	Close(ctx context.Context) (err error)
}

type Querier interface {
	// L2 cache backend
	Set(ctx context.Context, key, value string) (err error)
	Get(ctx context.Context, key string) (value string, err error)
	Del(ctx context.Context, key string) (err error)
	GetAll(ctx context.Context) (values map[string]string, err error)

	GetUser(ctx context.Context, id string) (*models.User, error)
	SetUser(ctx context.Context, userID string) error
	DelUser(ctx context.Context, id string) error
}

func UpsertUser(ctx context.Context, d Database, userID string) error {
	if _, err := d.Querier().GetUser(ctx, userID); err == nil {
		return nil
	}

	if err := d.Querier().SetUser(ctx, userID); err != nil {
		return err
	}

	return nil
}
