package sessions

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims[ClaimsData any] struct {
	SessionID   string
	SessionData ClaimsData
	jwt.RegisteredClaims
}

func (s *SessionStore[ClaimsData]) newJwt(sessionID string, sessionData ClaimsData) (string, error) {
	method := jwt.SigningMethodHS256

	now := time.Now()
	expires := now.Add(s.params.SessionTTL)

	claims := &Claims[ClaimsData]{
		SessionID:   sessionID,
		SessionData: sessionData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	key := []byte(s.params.StoreSecret)

	return jwt.NewWithClaims(method, claims).SignedString(key)
}

func (s *SessionStore[ClaimsData]) validateJwt(tokenString string) (*Claims[ClaimsData], error) {
	var empty *Claims[ClaimsData]
	claims := &Claims[ClaimsData]{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.params.StoreSecret), nil
	})
	if err != nil {
		return empty, err
	}

	claims, ok := token.Claims.(*Claims[ClaimsData])
	if !ok || !token.Valid {
		return empty, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
