package config

import (
	"app/sessions"
	"app/storage/kv"
)

type SessionData struct {
	OS   string
	Name string
}

func InitSessions(store kv.Store[sessions.Session]) *sessions.Store[SessionData] {
	client, err := sessions.NewStore[SessionData](sessions.StoreParams{
		Store:       store,
		StoreSecret: Config.App.Secret,
	}, true)
	if err != nil {
		panic("cannot start sessions client")
	}

	return client
}
