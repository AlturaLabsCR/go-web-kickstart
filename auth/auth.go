// Package auth
package auth

import "net/http"

// UserIDProvider returns a unique userID for an auth provider
type UserIDProvider interface {
	UserID(r *http.Request) (string, error)
}
