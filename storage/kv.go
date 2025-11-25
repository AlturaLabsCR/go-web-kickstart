package storage

import "sync"

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const ErrKeyNotFound = errStr("key not found")

type KVStorage[T any] interface {
	Set(key string, value T)
	Get(key string) (T, error)
	Delete(key string)
}

type KVMemory[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

func NewKVMemoryStore[T any]() *KVMemory[T] {
	return &KVMemory[T]{
		data: map[string]T{},
	}
}

func (m *KVMemory[T]) Set(key string, value T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *KVMemory[T]) Get(key string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	if !ok {
		var t T
		return t, ErrKeyNotFound
	}
	return value, nil
}

func (m *KVMemory[T]) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}
