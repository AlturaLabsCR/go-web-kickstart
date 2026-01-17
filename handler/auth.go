package handler

import (
	"context"
	"net/http"

	providers "app/auth"
	"app/config"
	"app/config/routes"
	"app/database"
	"app/database/models"
	"app/i18n"
	"app/templates/auth"
	"app/templates/base"

	"github.com/mileusna/useragent"
)

type SessionDataKey string

const SessionData SessionDataKey = "session.data"

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Sess().Validate(w, r); err == nil {
		http.Redirect(w, r, routes.Map[routes.Protected], http.StatusSeeOther)
		return
	} else {
		h.Log().Debug("error validating", "error", err)
	}

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

	ctx := r.Context()
	session, err := h.UpsertUser(ctx, h.DB(), r, userID)
	if err != nil {
		h.Log().Error("error upserting user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.Sess().Set(ctx, w, userID, session); err != nil {
		h.Log().Error("error setting session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log().Debug(
		"logged user in",
		"userID", userID,
		"sessionData", session,
	)

	http.Redirect(w, r, routes.Map[routes.Protected], http.StatusSeeOther)
}

func (h *Handler) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := h.Sess().Validate(w, r)
		if err != nil {
			h.Log().Debug("failed validation, redirecting", "error", err)
			http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), SessionData, sessionData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func (h *Handler) UpsertUser(ctx context.Context, d database.Database, r *http.Request, userID string) (*config.SessionData, error) {
	perms := []string{}
	var err error = nil
	userName := ""

	err = d.WithTx(ctx, func(q database.Querier) error {
		if _, err = q.GetUser(ctx, userID); err != nil {
			if err := q.SetUser(ctx, userID); err != nil {
				return err
			}

			initialized, err := q.GetConfig(ctx, models.ConfigInitialized)
			if err != nil {
				return err
			}

			if initialized == models.ConfigInitializedTrue {
				if err := q.SetRole(ctx, userID, "role.default"); err != nil {
					return err
				}
			} else {
				if err := q.SetRole(ctx, userID, "role.admin"); err != nil {
					return err
				}

				if err := q.SetConfig(ctx, models.ConfigInitialized, models.ConfigInitializedTrue); err != nil {
					return err
				}
			}
		}

		perms, err = q.GetPermissions(ctx, userID)
		if err != nil {
			return err
		}

		userName, err = q.GetUserName(ctx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	ua := useragent.Parse(r.UserAgent())
	session := &config.SessionData{
		UserID:   userID,
		UserName: userName,
		Agent:    ua.OS,
		Perms:    perms,
	}

	return session, nil
}
