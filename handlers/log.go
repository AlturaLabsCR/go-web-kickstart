package handlers

import (
	"net/http"
	"time"
)

const msgEndpointHit = "endpoint hit"

func (h *Handler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		h.Log().Debug(msgEndpointHit,
			"start", start,
			"end", end,
			"pattern", r.Pattern,
		)
	})
}
