package handler

import (
	"net/http"

	"app/config/routes"
	"app/templates"
	"app/templates/base"
)

func (h *Handler) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Sess().Validate(w, r); err != nil {
		http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
		return
	}

	ctx := r.Context()

	main := templates.HomeMain()

	params := base.HeadParams{
		LoadJS:      true,
		RobotsIndex: false,
	}

	tr := h.Tr(r)

	if err := base.Page(params, tr, main, routes.Map[routes.Protected]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
