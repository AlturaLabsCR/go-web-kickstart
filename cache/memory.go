package cache

import (
	"context"
	"sync"
)

type MemoryStore struct {
	m sync.Map // map[string]string
}

var _ Cache = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{m: sync.Map{}}
}

func (s *MemoryStore) Set(_ context.Context, key, value string) error {
	s.m.Store(key, value)
	return nil
}

func (s *MemoryStore) Get(_ context.Context, key string) (string, error) {
	v, ok := s.m.Load(key)
	if !ok {
		return "", ErrNotFound
	}
	return v.(string), nil
}

func (s *MemoryStore) Del(_ context.Context, key string) error {
	s.m.Delete(key)
	return nil
}

func (s *MemoryStore) GetAll(_ context.Context) (map[string]string, error) {
	values := make(map[string]string)

	s.m.Range(func(k, v any) bool {
		values[k.(string)] = v.(string)
		return true
	})

	return values, nil
}
