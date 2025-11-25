// Package kv implements key-value storage abstractions
package kv

type Store[T any] interface {
	Set(key string, value T) error
	Get(key string) (T, bool)
	Delete(key string) error
}
