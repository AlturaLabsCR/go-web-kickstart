// Package sessions
package sessions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"app/cache"
)

const (
	AccessTokenKey = "access_token"
	CSRFTokenKey   = "csrf_token"
	CSRFHeaderKey  = "X-CSRF-Token"

	ErrBadCache = errStr("invalid cache")
)

type Session struct {
	SessionUser string
	CSRFToken   string
	CreatedAt   int64
	LastUsedAt  int64
}

type StoreParams struct {
	Cache   cache.Cache
	L2Cache cache.Cache // optional

	NamespacePrefix string        // defaults to 'session:'
	CookiePrefix    string        // defaults to 'session.'
	CookiePath      string        // defaults to '/'
	SessionTTL      time.Duration // defaults to time.Hour * 24 * 30 = 1m
	CookieSameSite  http.SameSite // defaults tu http.SameSiteLaxMode
	StoreSecret     string        // defaults to generated string
}

type Store[T any] struct {
	params StoreParams
}

func NewStore[T any](ctx context.Context, params StoreParams) (*Store[T], error) {
	if params.Cache == nil {
		return nil, ErrBadCache
	}

	if params.L2Cache == nil {
		params.L2Cache = cache.NoopCache{}
	} else {
		elems, err := params.L2Cache.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		for key, value := range elems {
			err := params.Cache.Set(ctx, key, value)
			if err != nil {
				return nil, err
			}
		}
	}

	if params.NamespacePrefix == "" {
		params.NamespacePrefix = "session:"
	}

	if params.CookiePrefix == "" {
		params.CookiePrefix = "session."
	}

	if params.CookiePath == "" {
		params.CookiePath = "/"
	}

	if params.SessionTTL == 0 {
		params.SessionTTL = time.Hour * 24 * 30
	}

	if params.CookieSameSite == 0 {
		params.CookieSameSite = http.SameSiteLaxMode
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

func (s *Store[T]) Set(ctx context.Context, w http.ResponseWriter, sessionUser string, sessionData *T) error {
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

	now := time.Now().Unix()

	session := &Session{
		SessionUser: sessionUser,
		CSRFToken:   csrfToken,
		CreatedAt:   now,
		LastUsedAt:  now,
	}

	sessionStr, err := sessionToString(session)
	if err != nil {
		return err
	}

	_ = s.params.L2Cache.Set(ctx, s.params.NamespacePrefix+sessionID, sessionStr)

	return s.params.Cache.Set(ctx, s.params.NamespacePrefix+sessionID, sessionStr)
}

func (s *Store[T]) Validate(w http.ResponseWriter, r *http.Request) (*T, error) {
	ctx := r.Context()

	cookie, err := r.Cookie(s.params.CookiePrefix + AccessTokenKey)
	if err != nil {
		return nil, err
	}

	claims, err := s.validateJwt(cookie.Value)
	if err != nil {
		return nil, err
	}

	sessionStr, err := s.params.Cache.Get(ctx, s.params.NamespacePrefix+claims.SessionID)
	if err != nil {
		return nil, err
	}

	session, err := stringToSession(sessionStr)
	if err != nil {
		return nil, err
	}

	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		if t := r.Header.Get(CSRFHeaderKey); t != session.CSRFToken {
			return nil, fmt.Errorf("invalid CSRF token: expected %s, got %s", session.CSRFToken, t)
		}
	}

	if err := s.refresh(w, r, claims, session); err != nil {
		return nil, fmt.Errorf("failed to refresh session")
	}

	return &claims.SessionData, nil
}

func (s *Store[T]) Revoke(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(s.params.CookiePrefix + AccessTokenKey)
	if err != nil {
		return err
	}

	claims, err := s.validateJwt(cookie.Value)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.params.CookiePrefix + AccessTokenKey,
		Path:     s.params.CookiePath,
		SameSite: s.params.CookieSameSite,
		Expires:  time.Time{},
		Value:    "",
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     s.params.CookiePrefix + CSRFTokenKey,
		Path:     s.params.CookiePath,
		SameSite: s.params.CookieSameSite,
		Expires:  time.Time{},
		Value:    "",
		HttpOnly: false,
		Secure:   true,
	})

	if err := s.params.L2Cache.Del(r.Context(), s.params.NamespacePrefix+claims.SessionID); err != nil {
		return err
	}

	return s.params.Cache.Del(r.Context(), s.params.NamespacePrefix+claims.SessionID)
}
