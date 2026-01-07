package queries

import (
	"context"

	"app/database/sqlite/db"
)

func (sq *SqliteQuerier) Set(ctx context.Context, key, value string) error {
	return sq.queries.SetCache(ctx, db.SetCacheParams{
		Key:   key,
		Value: value,
	})
}

func (sq *SqliteQuerier) Get(ctx context.Context, key string) (string, error) {
	return sq.queries.GetCache(ctx, key)
}

func (sq *SqliteQuerier) Del(ctx context.Context, key string) error {
	return sq.queries.DelCache(ctx, key)
}

func (sq *SqliteQuerier) GetAll(ctx context.Context) (map[string]string, error) {
	elems, err := sq.queries.GetAllCache(ctx)
	if err != nil {
		return nil, err
	}

	all := map[string]string{}

	for _, elem := range elems {
		all[elem.Key] = elem.Value
	}

	return all, nil
}
