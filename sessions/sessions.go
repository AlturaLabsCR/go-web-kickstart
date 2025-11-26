// Package sessions handles user authorization and authentication
package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"app/storage/kv"
)

const (
	AccessTokenKey = "access_token"
	CSRFTokenKey   = "csrf_token"
	CSRFHeaderKey  = "X-CSRF-Token"
)

type Session struct {
	SessionUser string
	CSRFToken   string
	CreatedAt   time.Time
	LastUsedAt  time.Time
}

type StoreParams struct {
	SessionTTL time.Duration // defaults to 24 * 30 * time.Hour

	CookiePath     string        // defaults to '/'
	CookiePrefix   string        // defaults to 'session.'
	CookieSameSite http.SameSite // defaults tu http.SameSiteLaxMode

	Store       kv.Store[Session] // defaults to memory store
	StoreSecret string            // defaults to generated string
}

type Store[T any] struct {
	params StoreParams
}

func NewStore[T any](params StoreParams, secure bool) (*Store[T], error) {
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

	if params.StoreSecret == "" {
		s, err := generateToken()
		if err != nil {
			return nil, err
		}
		params.StoreSecret = s
	}

	return &Store[T]{params: params}, nil
}

func (s *Store[T]) Set(ctx context.Context, w http.ResponseWriter, sessionUser string, sessionData T) error {
	sessionID, err := generateToken()
	if err != nil {
		return err
	}

	if _, err := s.refreshAccessTokenCookie(w, sessionID, sessionData); err != nil {
		return err
	}

	csrfToken, err := s.refreshCSRFCookie(w)
	if err != nil {
		return err
	}

	now := time.Now()

	return s.params.Store.Set(ctx, sessionID, Session{
		SessionUser: sessionUser,
		CSRFToken:   csrfToken,
		CreatedAt:   now,
		LastUsedAt:  now,
	})
}

func (s *Store[T]) refresh(w http.ResponseWriter, r *http.Request, claims *Claims[T], session Session) error {
	ctx := r.Context()

	if _, err := s.refreshAccessTokenCookie(w, claims.SessionID, claims.SessionData); err != nil {
		return err
	}

	csrfToken, err := s.refreshCSRFCookie(w)
	if err != nil {
		return err
	}

	return s.params.Store.Set(ctx, claims.SessionID, Session{
		SessionUser: session.SessionUser,
		CSRFToken:   csrfToken,
		CreatedAt:   session.CreatedAt,
		LastUsedAt:  time.Now(),
	})
}

func (s *Store[T]) Validate(w http.ResponseWriter, r *http.Request) (T, error) {
	ctx := r.Context()
	var empty T

	cookie, err := r.Cookie(s.params.CookiePrefix + AccessTokenKey)
	if err != nil {
		return empty, err
	}

	claims, err := s.validateJwt(cookie.Value)
	if err != nil {
		return empty, err
	}

	session, err := s.params.Store.Get(ctx, claims.SessionID)
	if err != nil {
		return empty, err
	}

	if r.Method == http.MethodGet {
		return claims.SessionData, nil
	}

	csrfToken := r.Header.Get(CSRFHeaderKey)
	if csrfToken != session.CSRFToken {
		return empty, fmt.Errorf("invalid CSRF token")
	}

	if err := s.refresh(w, r, claims, session); err != nil {
		return empty, fmt.Errorf("failed to refresh session")
	}

	return claims.SessionData, nil
}

func (s *Store[T]) refreshAccessTokenCookie(w http.ResponseWriter, sessionID string, sessionData T) (string, error) {
	accessToken, err := s.newJwt(sessionID, sessionData)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.params.CookiePrefix + AccessTokenKey,
		Path:     s.params.CookiePath,
		SameSite: s.params.CookieSameSite,
		Expires:  time.Now().Add(s.params.SessionTTL),
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
	})

	return accessToken, nil
}

func (s *Store[T]) refreshCSRFCookie(w http.ResponseWriter) (string, error) {
	csrfToken, err := generateToken()
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.params.CookiePrefix + CSRFTokenKey,
		Path:     s.params.CookiePath,
		SameSite: s.params.CookieSameSite,
		Expires:  time.Now().Add(s.params.SessionTTL),
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   true,
	})

	return csrfToken, nil
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
