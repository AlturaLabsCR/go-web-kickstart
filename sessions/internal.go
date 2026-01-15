package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

func (s *Store[T]) refresh(w http.ResponseWriter, r *http.Request, claims *Claims[T], session *Session) error {
	ctx := r.Context()

	if _, err := s.refreshAccessTokenCookie(
		w,
		claims.SessionID,
		&claims.SessionData,
	); err != nil {
		return err
	}

	csrfToken, err := s.refreshCSRFCookie(w)
	if err != nil {
		return err
	}

	session.CSRFToken = csrfToken
	session.LastUsedAt = time.Now().Unix()

	sessionStr, err := sessionToString(session)
	if err != nil {
		return err
	}

	_ = s.params.L2Cache.Set(ctx, claims.SessionID, sessionStr)

	return s.params.Cache.Set(ctx, claims.SessionID, sessionStr)
}

func (s *Store[T]) refreshAccessTokenCookie(w http.ResponseWriter, sessionID string, sessionData *T) (string, error) {
	accessToken, err := s.newJwt(sessionID, *sessionData)
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

func sessionToString(session *Session) (string, error) {
	s, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func stringToSession(sessionStr string) (*Session, error) {
	var s Session
	if err := json.Unmarshal([]byte(sessionStr), &s); err != nil {
		return nil, err
	}
	return &s, nil
}
