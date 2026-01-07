package handler

import (
	"net/http"

	"app/config/routes"
	"app/templates"
	"app/templates/base"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	main := templates.HomeMain()

	params := base.HeadParams{
		LoadJS:      true,
		RobotsIndex: true,
	}

	tr := h.Tr(r)

	if err := base.Page(params, tr, main, routes.Map[routes.Root]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}

func (h *Handler) AboutPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	main := templates.HomeMain()

	params := base.HeadParams{
		LoadJS:      true,
		RobotsIndex: true,
	}

	tr := h.Tr(r)

	if err := base.Page(params, tr, main, routes.Map[routes.About]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
