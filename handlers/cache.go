package handlers

import (
	"net/http"
)

func (h *Handler) CachePolicy(next http.Handler) http.HandlerFunc {
	const cacheHeader = "Cache-Control"
	cacheValue := "no-store"

	if h.Production() {
		cacheValue = "public, max-age=31536000, immutable"
	}

	h.Log().Debug("caching policy",
		"header", cacheHeader,
		"value", cacheValue,
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheHeader, cacheValue)
		next.ServeHTTP(w, r)
	})
}
