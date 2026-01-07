package handler

import (
	"net/http"

	"app/config/routes"
	"app/templates"
	"app/templates/base"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	main := templates.HomeMain()

	params := base.HeadParams{
		LoadJS:      true,
		RobotsIndex: true,
	}

	tr := h.Tr(r)

	if err := base.Page(params, tr, main, routes.Map[routes.Login]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
