// Package database
package database

import (
	"context"

	"app/database/models"
)

type Database interface {
	WithTx(ctx context.Context, fn func(q Querier) error) (err error)
	Exec(ctx context.Context, sql string) (err error)
	Close(ctx context.Context) (err error)
}

type Querier interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	SetUser(ctx context.Context, id string, user *models.User) error
	DelUser(ctx context.Context, id string) error
}
