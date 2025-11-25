// Package sessions handles user authorization and authentication
package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"app/storage/kv"
)

type StoreParams[T any] struct {
	SessionTTL time.Duration // defaults to time.Hour
	RefreshTTL time.Duration // defaults to 24 * 30 * time.Hour

	CookiePath     string        // defaults to '/'
	CookiePrefix   string        // defaults to 'session.'
	CookieSameSite http.SameSite // defaults tu http.SameSiteLaxMode

	Store       kv.Store[T] // defaults to memory store
	StoreSecret string      // defaults to auto-generated string
}

type Store[T any] struct {
	params StoreParams[T]
}

func NewStore[T any](params StoreParams[T]) (*Store[T], error) {
	if params.SessionTTL == 0 {
		params.SessionTTL = time.Hour
	}

	if params.RefreshTTL == 0 {
		params.RefreshTTL = 24 * 30 * time.Hour
	}

	if params.CookiePath == "" {
		params.CookiePath = "/"
	}

	if params.CookiePrefix == "" {
		params.CookiePrefix = "session."
	}

	if params.CookieSameSite == 0 {
		params.CookieSameSite = http.SameSiteLaxMode
	}

	if params.StoreSecret == "" {
		if secret, err := generateToken(); err == nil {
			params.StoreSecret = secret
		} else {
			return nil, err
		}
	}

	if params.Store == nil {
		params.Store = kv.NewMemoryStore[T]()
	}

	return &Store[T]{params: params}, nil
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
