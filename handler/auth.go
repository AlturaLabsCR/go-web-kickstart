package handler

import (
	"net/http"

	providers "app/auth"
	"app/config"
	"app/config/routes"
	"app/database"
	"app/i18n"
	"app/templates/auth"
	"app/templates/base"

	"github.com/mileusna/useragent"
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

func (h *Handler) LoginWithFacebook(w http.ResponseWriter, r *http.Request) {
	h.Log().Debug("authenticating with facebook")

	provider := &providers.FacebookProvider{
		AppID:      config.Config.AuthProviders.Facebook.AppID,
		AppSecret:  config.Config.AuthProviders.Facebook.AppSecret,
		HTTPClient: http.DefaultClient,
	}

	h.loginWithProvider(provider, w, r)
}

func (h *Handler) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	h.Log().Debug("authenticating with google")

	provider := &providers.GoogleProvider{
		ClientID: config.Config.AuthProviders.Google.ClientID,
	}

	h.loginWithProvider(provider, w, r)
}

func (h *Handler) loginWithProvider(provider providers.UserIDProvider, w http.ResponseWriter, r *http.Request) {
	if _, err := h.Sess().Validate(w, r); err == nil {
		http.Redirect(w, r, routes.Map[routes.Protected], http.StatusSeeOther)
		return
	}

	userID, err := provider.UserID(r)
	if err != nil {
		h.Log().Debug("error getting userID", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ua := useragent.Parse(r.UserAgent())
	sessionData := &config.SessionData{
		UserID: userID,
		Agent:  ua.OS,
		Perms:  []string{},
	}

	ctx := r.Context()
	if err := database.UpsertUser(ctx, h.DB(), userID); err != nil {
		h.Log().Debug("error upserting user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log().Debug(
		"logged user in",
		"userID", userID,
		"sessionData", sessionData,
	)

	http.Redirect(w, r, routes.Map[routes.Protected], http.StatusSeeOther)
}
