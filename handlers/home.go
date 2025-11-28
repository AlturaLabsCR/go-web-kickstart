package handlers

import (
	"net/http"

	"app/templates"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tr := h.Translator(r)

	content := templates.Home(tr)
	loadFrameworks := false

	templates.Base(content, loadFrameworks).Render(ctx, w)
}
