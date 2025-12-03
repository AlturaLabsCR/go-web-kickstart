package handlers

import (
	"net/http"

	"app/templates"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tr := h.Translator(r)

	content := templates.Home(tr)

	templates.Base(content, templates.BaseTemplateParams{
		Description: []templates.Tag{
			{Name: "", Value: tr("home.description")},
		},
		RobotsIndex: true,
	}).Render(ctx, w)
}
