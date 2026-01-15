package queries

import (
	"context"
	"database/sql"

	"app/database/models"
	"app/database/sqlite/db"
)

func (sq *SqliteQuerier) GetConfigs(ctx context.Context) ([]models.Config, error) {
	confs, err := sq.queries.GetConfigs(ctx)
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

func (sq *SqliteQuerier) GetConfig(ctx context.Context, name string) (string, error) {
	value, err := sq.cache.Get(ctx, name)
	if err == nil {
		return value, nil
	}

	value, err = sq.queries.GetConfig(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return value, nil
}

func (sq *SqliteQuerier) SetConfig(ctx context.Context, name, value string) error {
	if err := sq.cache.Set(ctx, name, value); err != nil {
		return err
	}

	return sq.queries.SetConfig(ctx, db.SetConfigParams{Name: name, Value: value})
}
