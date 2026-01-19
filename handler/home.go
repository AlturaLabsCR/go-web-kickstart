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

	head := base.HeadParams{
		LoadJS:      false,
		RobotsIndex: true,
	}

	page := base.PageParams{
		Head: head,
		Body: base.BodyParams{
			Content: main,
			Active:  routes.Map[routes.Root],
		},
	}

	tr := h.Tr(r)

	if err := base.Page(tr, page).Render(ctx, w); err != nil {
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

	head := base.HeadParams{
		Subtitle:    tr("nav.about"),
		LoadJS:      false,
		RobotsIndex: true,
	}

	page := base.PageParams{
		Head: head,
		Body: base.BodyParams{
			Content: main,
			Active:  routes.Map[routes.About],
		},
	}

	if err := base.Page(tr, page).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
