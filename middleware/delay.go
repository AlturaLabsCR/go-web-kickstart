package middleware

import (
	"net/http"
	"time"
)

func Delay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(90 * time.Millisecond):
		}
		next.ServeHTTP(w, r)
	})
}
