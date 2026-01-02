// Package database
package database

import "context"

type Database interface {
	WithTx(ctx context.Context, fn func(q Querier) error) error
	ExecSQL(context.Context, string) error
	Close(context.Context)
}

type Querier interface {
	internalQuerier
}
