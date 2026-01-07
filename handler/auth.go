package handler

import (
	"net/http"

	"app/config"
	"app/config/routes"
	"app/i18n"
	"app/templates/auth"
	"app/templates/base"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
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

	authParams := auth.LoginParams{
		GoogleClientID:       config.Config.AuthProviders.Google.ClientID,
		GoogleVerifyEndpoint: routes.Map[routes.GoogleAuth],
		FacebookAuthParams: auth.FacebookAuthParams{
			AppID:    config.Config.AuthProviders.Facebook.AppID,
			Locale:   locale,
			Version:  config.FacebookAPIVersion,
			Endpoint: routes.Map[routes.FacebookAuth],
		},
	}

	main := auth.LoginMain(authParams)

	params := base.HeadParams{
		LoadJS:      true,
		RobotsIndex: true,
	}

	tr := h.Tr(r)

	if err := base.Page(params, tr, main, routes.Map[routes.Login]).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
