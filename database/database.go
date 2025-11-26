// Package database abstracts different engines into one Database interface
package database

import (
	"context"
)

type Database interface {
	Close(context.Context)
	upsertUser(context.Context, string) error
}

type errStr string

func (e errStr) Error() string {
	return string(e)
}

func UpsertUser(db Database, ctx context.Context, userID string) error {
	return db.upsertUser(ctx, userID)
}
