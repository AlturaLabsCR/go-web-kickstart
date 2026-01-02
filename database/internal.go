package database

import (
	"context"

	"app/cache"
)

type NoopCache struct{}

func (NoopCache) Get(ctx context.Context, key string) (string, error) {
	return "", cache.ErrNotFound
}

func (NoopCache) Set(ctx context.Context, key string, value string) error {
	return nil
}

func (NoopCache) Del(ctx context.Context, key string) error {
	return nil
}

func (NoopCache) GetAll(ctx context.Context) (map[string]string, error) {
	return nil, cache.ErrNotFound
}
