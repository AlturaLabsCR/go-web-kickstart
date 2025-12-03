// Package utils has reusable helper methods
package utils

import (
	"fmt"
	"io"
	"net/http"
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

func InspectReader(r io.Reader, maxSize int64) (mime string, size int64, data []byte, err error) {
	if maxSize <= 0 {
		return "", 0, nil, fmt.Errorf("maxSize must be > 0")
	}

	limited := &io.LimitedReader{
		R: r,
		N: maxSize + 1,
	}

	data, err = io.ReadAll(limited)
	if err != nil {
		return "", 0, nil, err
	}

	size = int64(len(data))

	mime = http.DetectContentType(data)

	return mime, size, data, nil
}
