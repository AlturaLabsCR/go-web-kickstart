// Package database abstracts different engines into one Database interface
package database

import (
	"context"
	"slices"

	"app/utils"
)

type Database interface {
	Close(context.Context)
	insertUser(context.Context, string) (int64, error)
	selectUserEmails(context.Context) ([]string, error)
}

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const (
	ErrDuplicateEmail = errStr("duplicate email")
)

func InsertUser(db Database, ctx context.Context, userEmail string) (int64, error) {
	userEmail, err := utils.ParseEmail(userEmail)
	if err != nil {
		return 0, err
	}

	// TODO: Cache this
	if registeredEmails, err := db.selectUserEmails(ctx); err != nil {
		return 0, err
	} else {
		if slices.Contains(registeredEmails, userEmail) {
			return 0, ErrDuplicateEmail
		}
	}

	return db.insertUser(ctx, userEmail)
}
