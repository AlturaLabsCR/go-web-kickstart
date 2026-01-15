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

	// Configs
	GetConfigs(ctx context.Context) ([]models.Config, error)
	GetConfig(ctx context.Context, name string) (value string, err error)
	SetConfig(ctx context.Context, name, value string) error

	// TODO: Cache permissions
	SetRole(ctx context.Context, userID, roleName string) (err error)
	GetPermissions(ctx context.Context, userID string) (perms []string, err error)
}

func UpsertUser(ctx context.Context, d Database, userID string) (perms []string, err error) {
	err = d.WithTx(ctx, func(q Querier) error {
		_, err := q.GetUser(ctx, userID)
		if err != nil {
			if err := q.SetUser(ctx, userID); err != nil {
				return err
			}

			initialized, err := q.GetConfig(ctx, "config.initialized")
			if err != nil {
				return err
			}

			if initialized == "true" {
				if err := q.SetRole(ctx, userID, "role.default"); err != nil {
					return err
				}
			} else {
				if err := q.SetRole(ctx, userID, "role.admin"); err != nil {
					return err
				}

				if err := q.SetConfig(ctx, "config.initialized", "true"); err != nil {
					return err
				}
			}
		}

		perms, err = q.GetPermissions(ctx, userID)
		return err
	})

	return perms, err
}
