// Package kv implements key-value storage abstractions
package kv

import "context"

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const ErrNotFound = errStr("not found")

type Store[T any] interface {
	Set(ctx context.Context, key string, v T) error
	Get(ctx context.Context, key string) (T, error)
	Delete(ctx context.Context, key string) error
	GetElems(ctx context.Context) (map[string]T, error)
}
