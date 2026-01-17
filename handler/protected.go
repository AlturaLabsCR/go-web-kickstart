package handler

import (
	"net/http"

	"app/config"
	"app/config/routes"
	"app/templates/base"
	"app/templates/protected"
)

func (h *Handler) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, ok := r.Context().Value(SessionData).(*config.SessionData)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	tr := h.Tr(r)

	main := protected.ProtectedMain(tr, session, r.URL.Path)

	params := base.HeadParams{
		Subtitle:    tr("nav.account"),
		LoadJS:      true,
		RobotsIndex: false,
	}

	if err := base.Page(params, tr, main, routes.Map[routes.Protected]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
