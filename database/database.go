// Package database abstracts different engines into one Database interface
package database

import (
	"context"
)

type Database interface {
	ExecSQL(context.Context, string) error
	Close(context.Context)
	upsertUser(context.Context, string) error
}

func UpsertUser(db Database, ctx context.Context, userID string) error {
	return db.upsertUser(ctx, userID)
}
