// Package database abstracts different engines into one Database interface
package database

import (
	"context"
	"slices"

	"app/utils"
)

type Database interface {
	Close(context.Context)
	insertOwner(context.Context, string) (int64, error)
	selectOwnerEmails(context.Context) ([]string, error)
}

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const (
	ErrDuplicateEmail = errStr("duplicate email")
)

func InsertOwner(db Database, ctx context.Context, ownerEmail string) (int64, error) {
	ownerEmail, err := utils.ParseEmail(ownerEmail)
	if err != nil {
		return 0, err
	}

	// TODO: Cache this
	if registeredEmails, err := db.selectOwnerEmails(ctx); err != nil {
		return 0, err
	} else {
		if slices.Contains(registeredEmails, ownerEmail) {
			return 0, ErrDuplicateEmail
		}
	}

	return db.insertOwner(ctx, ownerEmail)
}
