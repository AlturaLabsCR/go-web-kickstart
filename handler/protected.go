package handler

import (
	"net/http"

	"app/config/routes"
	"app/templates/base"
	"app/templates/protected"
)

func (h *Handler) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Sess().Validate(w, r); err != nil {
		http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
		return
	}

	ctx := r.Context()

	main := protected.ProtectedMain()

	tr := h.Tr(r)

	params := base.HeadParams{
		Subtitle:    tr("nav.account"),
		LoadJS:      true,
		RobotsIndex: false,
	}

	if err := base.Page(params, tr, main, routes.Map[routes.Protected]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
