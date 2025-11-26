// Package kv implements key-value storage abstractions
package kv

import "context"

type Store[T any] interface {
	Set(context.Context, string, T) error
	Get(context.Context, string) (T, error)
	Delete(context.Context, string) error
}
