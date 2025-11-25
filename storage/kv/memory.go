package kv

import "sync"

type MemoryStore[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		data: map[string]T{},
	}
}

func (m *MemoryStore[T]) Set(key string, value T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *MemoryStore[T]) Get(key string) (T, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	return value, ok
}

func (m *MemoryStore[T]) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}
