// Package sessions handles user authorization and authentication
package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"app/storage/kv"
)

const (
	AccessTokenKey = "access_token"
	CSRFTokenKey   = "csrf_token"
)

type Session struct {
	AccessToken string
	CSRFToken   string
	ExpiresAt   time.Time
}

type StoreParams struct {
	SessionTTL time.Duration // defaults to time.Hour

	CookiePath     string        // defaults to '/'
	CookiePrefix   string        // defaults to 'session.'
	CookieSameSite http.SameSite // defaults tu http.SameSiteLaxMode

	Store       kv.Store[Session] // defaults to memory store
	StoreSecret string            // defaults to generated string
}

type Store struct {
	params StoreParams
}

func NewStore(params StoreParams, secure bool) (*Store, error) {
	if params.SessionTTL == 0 {
		params.SessionTTL = 24 * 30 * time.Hour
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

	if params.Store == nil {
		params.Store = kv.NewMemoryStore[Session]()
	}

	return &Store{params: params}, nil
}

func (s *Store) Set(w http.ResponseWriter, sessionID string) error {
	accessToken := "my-signed-jwt-token"

	csrfToken, err := generateToken()
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(s.params.SessionTTL)

	http.SetCookie(w, &http.Cookie{
		Name:     s.params.CookiePrefix + AccessTokenKey,
		Path:     s.params.CookiePath,
		SameSite: s.params.CookieSameSite,
		Expires:  expiresAt,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     s.params.CookiePrefix + CSRFTokenKey,
		Path:     s.params.CookiePath,
		SameSite: s.params.CookieSameSite,
		Expires:  expiresAt,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   true,
	})

	return s.params.Store.Set(sessionID, Session{
		AccessToken: accessToken,
		CSRFToken:   csrfToken,
		ExpiresAt:   expiresAt,
	})
}

func (s *Store) Validate(r *http.Request, enforceCSRFProtection bool) error {
	// 1. Get sessionID from JWT claims
	// If err != nil, sessionID is known to be valid as the claims are
	// signed by a server secret, so:

	// 2. Check for enforceCSRFProtection in the *header*

	// 3. If the SessionTTL has not passed, re-roll Session, else error

	return nil
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
