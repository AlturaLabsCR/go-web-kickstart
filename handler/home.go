package handler

import (
	"net/http"

	"app/config"
	"app/config/routes"
	"app/i18n"
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

	locale := ""

	for _, lang := range i18n.RequestLanguages(r) {
		if loc, ok := config.SupportedLocales[lang.Tag]; ok {
			locale = loc
			break
		}
	}

	if locale == "" {
		locale = config.SupportedLocales[config.DefaultLocale]
	}

	tr := h.Tr(r)

	main := templates.AboutMain(tr, locale)

	params := base.HeadParams{
		LoadJS:      true,
		RobotsIndex: true,
	}

	if err := base.Page(params, tr, main, routes.Map[routes.About]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
