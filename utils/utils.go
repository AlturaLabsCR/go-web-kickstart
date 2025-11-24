// Package utils has reusable helper methods
package utils

import (
	"net/mail"
	"strings"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const ErrBadEmail = errStr("failed to parse email")

func ParseEmail(email string) (string, error) {
	email = strings.TrimSpace(email)

	if len(email) < 5 || len(email) > 64 {
		return "", ErrBadEmail
	}

	email = strings.ToLower(email)

	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrBadEmail
	}

	return email, nil
}
