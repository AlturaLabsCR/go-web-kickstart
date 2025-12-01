package handlers

import (
	"net/http"

	"app/config"
	"app/templates"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := h.Translator(r)

	// if the session is valid, redirect
	if _, err := h.Sessions().Validate(w, r); err == nil {
		templates.Redirect(config.Endpoints[config.ProtectedPath]).Render(ctx, w)
		return
	}

	params := templates.LoginParams{
		GoogleClientID:         config.Environment[config.EnvGoogleClientID],
		GoogleVerifyEndpoint:   config.Endpoints[config.AuthWithGooglePath],
		FacebookAppID:          config.Environment[config.EnvFacebookAppID],
		FacebookVerifyEndpoint: config.Endpoints[config.AuthWithFacebookPath],
	}

	content := templates.Login(tr, params)
	loadFrameworks := false

	templates.Base(content, loadFrameworks).Render(ctx, w)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions().Revoke(w, r)

	templates.Redirect(config.Endpoints[config.LoginPath]).Render(r.Context(), w)
}
