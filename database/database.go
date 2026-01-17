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

	// Users
	GetUser(ctx context.Context, id string) (*models.User, error)
	SetUser(ctx context.Context, userID string) error
	DelUser(ctx context.Context, id string) error
	UpsertUserName(ctx context.Context, userName, userID string) error
	GetUserMeta(ctx context.Context, userID string) (meta *models.UserMeta, err error)

	// Configs
	GetConfigs(ctx context.Context) ([]models.Config, error)
	GetConfig(ctx context.Context, name string) (value string, err error)
	SetConfig(ctx context.Context, name, value string) error

	// TODO: Cache permissions
	SetRole(ctx context.Context, userID, roleName string) (err error)
	GetPermissions(ctx context.Context, userID string) (perms []string, err error)
}
