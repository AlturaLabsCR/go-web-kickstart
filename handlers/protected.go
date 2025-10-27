package handlers

import (
	"context"
	"net/http"
)

type ctxEmail struct{}

func (h *Handler) Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, err := h.Sessions.JWTValidate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			ctxEmail{},
			email,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
