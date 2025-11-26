package config

import (
	"app/sessions"
	"app/storage/kv"
)

type SessionData struct {
	OS string
}

func InitSessions(store kv.Store[sessions.Session]) *sessions.Store[SessionData] {
	client, err := sessions.NewStore[SessionData](sessions.StoreParams{
		Store:       store,
		StoreSecret: Environment[EnvSecret],
	}, Environment[EnvProd] == "1")
	if err != nil {
		panic("cannot start sessions client")
	}

	return client
}
