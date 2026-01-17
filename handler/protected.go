package handler

import (
	"net/http"

	"app/config/routes"
	"app/templates/base"
	"app/templates/protected"
)

func (h *Handler) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionData, ok := h.Sess().Data(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	sessionAttrs, ok := h.Sess().Attrs(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	tr := h.Tr(r)

	userMeta, err := h.DB().Querier().GetUserMeta(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Error(w, "error getting user meta", http.StatusInternalServerError)
		return
	}

	main := protected.ProtectedMain(tr, protected.ProtectedParams{
		User:  userMeta,
		Attrs: sessionAttrs,
		Data:  sessionData,
	}, r.URL.Path)

	params := base.HeadParams{
		Subtitle:    tr("nav.account"),
		LoadJS:      true,
		RobotsIndex: false,
	}

	if err := base.Page(params, tr, main, routes.Map[routes.Protected]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
