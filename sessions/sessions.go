// Package sessions handles user authorization and authentication
package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"app/storage"
)

type session struct {
	accessToken, csrfToken string
}

type StoreParams struct {
	SessionTTL time.Duration // defaults to time.Hour
	RefreshTTL time.Duration // defaults to 24 * 30 * time.Hour

	CookiePath     string        // defaults to '/'
	CookiePrefix   string        // defaults to 'session.'
	CookieSameSite http.SameSite // defaults tu http.SameSiteLaxMode

	StoreSecret string // defaults to auto-generated string
}

type Store[T any] struct {
	params  StoreParams
	storage storage.KVStorage[session] // defaults to storage.KVMemory[Session]
}

func NewStore[T any](params StoreParams) (*Store[T], error) {
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

	storage := storage.NewKVMemoryStore[session]()

	return &Store[T]{params: params, storage: storage}, nil
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
