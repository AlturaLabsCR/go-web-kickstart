// Package cache implements key-value-like interface 'Cache'
package cache

import (
	"context"
)

const (
	ErrNotFound = errStr("key not found")
)

type Cache interface {
	Set(ctx context.Context, key, value string) (err error)
	Get(ctx context.Context, key string) (value string, err error)
	Del(ctx context.Context, key string) (err error)
	GetAll(ctx context.Context) (values map[string]string, err error)
}
