package config

import (
	"app/cache"
	"app/database"
	"app/sessions"
)

func InitSessions(db database.Database) (*sessions.Store[SessionData], error) {
	var empty *sessions.Store[SessionData]

	params := sessions.StoreParams{
		Cache: cache.NewMemoryStore(),

		// TODO: allow db to be used as L2 cache
		// L2Cache cache.Cache // optional

		StoreSecret: Config.Sessions.Secret,
	}

	store, err := sessions.NewStore[SessionData](params)
	if err != nil {
		return empty, err
	}

	return store, nil
}
