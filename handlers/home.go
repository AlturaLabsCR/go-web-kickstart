package handlers

import (
	"net/http"

	"app/templates"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tr := h.Translator(r)

	content := templates.Home(tr)

	templates.Base(content).Render(ctx, w)
}
