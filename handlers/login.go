package handlers

import (
	"net/http"

	"app/auth"
	"app/config"
	"app/i18n"
	"app/templates"
)

const defaultLocale = "en"

var supportedLocales = map[string]string{
	"en": "en_US",
	"es": "es_LA",
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := h.Translator(r)

	// if the session is valid, redirect
	if _, err := h.Sessions().Validate(w, r); err == nil {
		templates.Redirect(config.Endpoints[config.ProtectedPath]).Render(ctx, w)
		return
	}

	locale := ""

	for _, lang := range i18n.RequestLanguages(r) {
		if loc, ok := supportedLocales[lang.Tag]; ok {
			locale = loc
			break
		}
	}

	if locale == "" {
		locale = supportedLocales[defaultLocale]
	}

	params := templates.LoginParams{
		GoogleClientID:       config.Config.AuthProviders.Google.ClientID,
		GoogleVerifyEndpoint: config.Endpoints[config.AuthWithGooglePath],
		FacebookAuthParams: templates.FacebookAuthParams{
			AppID:    config.Config.AuthProviders.Facebook.AppID,
			Locale:   locale,
			Version:  auth.FacebookAPIVersion,
			Endpoint: config.Endpoints[config.AuthWithFacebookPath],
		},
	}

	content := templates.Login(tr, params)
	loadFrameworks := false

	templates.Base(content, loadFrameworks).Render(ctx, w)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions().Revoke(w, r)

	templates.Redirect(config.Endpoints[config.LoginPath]).Render(r.Context(), w)
}
