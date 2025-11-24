// Package database abstracts different engines into one Database interface
package database

import "context"

type Database interface {
	InsertOwner(context.Context, string) (int64, error)
}
