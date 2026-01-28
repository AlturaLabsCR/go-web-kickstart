package config

import (
	"context"

	"app/cache"
	"app/database"
	"app/sessions"
)

const ProductionEnv = "production"

func InitSessions(ctx context.Context, db database.Database) (*sessions.Store[SessionData], error) {
	var empty *sessions.Store[SessionData]

	params := sessions.StoreParams{
		Cache:        cache.NewMemoryStore(),
		L2Cache:      db.Querier(),
		StoreSecret:  Config.Sessions.Secret,
		SecureCookie: Config.App.Env == ProductionEnv,
	}

	store, err := sessions.NewStore[SessionData](ctx, params)
	if err != nil {
		return empty, err
	}

	return store, nil
}
