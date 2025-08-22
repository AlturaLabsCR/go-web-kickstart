package handlers

import (
	"net/http"

	"app/templates"
)

func (h *Handler) RenderName(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("endpoint hit", "path", r.URL.Path)

	ctx := r.Context()

	name := r.PathValue("name")

	if err := templates.Hello(name).Render(ctx, w); err != nil {
		h.Log().Info("template render", "error", err)
	}
}
