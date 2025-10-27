package handlers

import (
	"net/http"

	"app/templates"
)

func (h *Handler) RenderName(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("endpoint hit", "path", r.URL.Path)

	ctx := r.Context()

	name := r.PathValue("name")

	content := templates.Hello(name)

	templates.Base(content).Render(ctx, w)
}
