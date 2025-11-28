// Package auth has common methods for authenticating users
package auth

import (
	"fmt"
	"net/http"

	"google.golang.org/api/idtoken"
)

func GetGoogleID(r *http.Request, clientID string) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}

	token := r.FormValue("token")
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	payload, err := idtoken.Validate(
		r.Context(),
		token,
		clientID,
	)
	if err != nil {
		return "", err
	}

	return payload.Subject, nil
}
