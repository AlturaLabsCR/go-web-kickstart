package queries

import (
	"context"

	"app/database/postgres/db"
)

func (pq *PostgresQuerier) Set(ctx context.Context, key, value string) error {
	return pq.queries.SetCache(ctx, db.SetCacheParams{
		Key:   key,
		Value: value,
	})
}

func (pq *PostgresQuerier) Get(ctx context.Context, key string) (string, error) {
	return pq.queries.GetCache(ctx, key)
}

func (pq *PostgresQuerier) Del(ctx context.Context, key string) error {
	return pq.queries.DelCache(ctx, key)
}

func (pq *PostgresQuerier) GetAll(ctx context.Context) (map[string]string, error) {
	elems, err := pq.queries.GetAllCache(ctx)
	if err != nil {
		return nil, err
	}

	all := map[string]string{}

	for _, elem := range elems {
		all[elem.Key] = elem.Value
	}

	return all, nil
}
