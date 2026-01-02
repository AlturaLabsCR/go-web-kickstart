package auth

import (
	"fmt"
	"net/http"

	"google.golang.org/api/idtoken"
)

const googlePrefix = "g:"

type GoogleProvider struct {
	ClientID string
}

func (p *GoogleProvider) UserID(r *http.Request) (string, error) {
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
		p.ClientID,
	)
	if err != nil {
		return "", err
	}

	return googlePrefix + payload.Subject, nil
}
