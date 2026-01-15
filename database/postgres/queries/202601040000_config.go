package queries

import (
	"context"

	"app/database/models"
	"app/database/postgres/db"

	"github.com/jackc/pgx/v5"
)

func (pq *PostgresQuerier) GetConfigs(ctx context.Context) ([]models.Config, error) {
	confs, err := pq.queries.GetConfigs(ctx)
	if err != nil {
		return nil, err
	}

	out := []models.Config{}
	for _, conf := range confs {
		this := models.Config{
			Name:  conf.Name,
			Value: conf.Value,
		}
		out = append(out, this)
	}

	return out, nil
}

func (pq *PostgresQuerier) GetConfig(ctx context.Context, name string) (string, error) {
	value, err := pq.cache.Get(ctx, models.ConfigCacheScopePrefix+name)
	if err == nil {
		return value, nil
	}

	value, err = pq.queries.GetConfig(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return value, nil
}

func (pq *PostgresQuerier) SetConfig(ctx context.Context, name, value string) error {
	if err := pq.cache.Set(ctx, models.ConfigCacheScopePrefix+name, value); err != nil {
		return err
	}

	return pq.queries.SetConfig(ctx, db.SetConfigParams{Name: name, Value: value})
}
