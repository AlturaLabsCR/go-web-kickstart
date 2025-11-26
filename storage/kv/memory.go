package kv

import (
	"context"
	"sync"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const ErrNotFound = errStr("not found")

type MemoryStore[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		data: map[string]T{},
	}
}

func (m *MemoryStore[T]) Set(_ context.Context, key string, value T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *MemoryStore[T]) Get(_ context.Context, key string) (T, error) {
	var empty T

	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	if !ok {
		return empty, ErrNotFound
	}

	return value, nil
}

func (m *MemoryStore[T]) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}
