package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type CachePolicy struct {
	Enabled bool
	Public  bool
	MaxAge  int64
}

func Cache(policy CachePolicy) func(http.Handler) http.Handler {
	value := "no-store"
	if policy.Enabled {
		parts := []string{}

		parts = append(parts, "immutable")

		if policy.Public {
			parts = append(parts, "public")
		} else {
			parts = append(parts, "private")
		}

		maxAge := policy.MaxAge
		if maxAge <= 0 {
			maxAge = int64(365 * 24 * 60 * 60)
		}

		parts = append(parts, fmt.Sprintf("max-age=%d", maxAge))

		value = strings.Join(parts, ", ")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)
			next.ServeHTTP(w, r)
		})
	}
}
