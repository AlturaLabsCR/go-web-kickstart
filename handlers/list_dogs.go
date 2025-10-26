package handlers

import (
	"net/http"

	"app/internal/db"
	"app/templates"
)

func (h *Handler) ListDogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queries := db.New(h.DB())

	dogs, err := queries.AllDogs(ctx)
	if err != nil {
		h.Log().Error("list dogs", "error", err)
	}

	if err := templates.Dogs(dogs).Render(ctx, w); err != nil {
		h.Log().Error("list dogs", "error", err)
	}
}
