package handlers

import (
	"context"
	"net/http"

	"app/config"
	"app/templates"
)

type handlerCtxKey string

const sessionDataCtxKey = handlerCtxKey("session.data")

func (h *Handler) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := h.Translator(r)

	content := templates.Protected(tr)

	templates.Base(content).Render(ctx, w)
}

func (h *Handler) Validate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := h.Sessions().Validate(w, r)
		if err != nil {
			h.Log().Debug("error validating session", "error", err)
			templates.Redirect(config.Endpoints[config.LoginPath]).Render(ctx, w)
			return
		}

		h.Log().Debug("validated session")

		ctx = context.WithValue(ctx, sessionDataCtxKey, data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
