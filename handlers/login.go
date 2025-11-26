package handlers

import (
	"net/http"

	"app/config"
	"app/database"
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

	content := templates.Login(tr)

	templates.Base(content).Render(ctx, w)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := h.Translator(r)

	// if the session is valid, redirect
	if _, err := h.Sessions().Validate(w, r); err == nil {
		templates.Redirect(config.Endpoints[config.ProtectedPath]).Render(ctx, w)
		return
	}

	// TODO: Validate user with oauth provider
	sessionUser := "sample_user_id"

	// TODO: Get session data from client (headers, request, body, etc)
	sessionData := config.SessionData{
		OS:       "linux",
		Location: "New york",
	}

	// from now on the user is authenticated

	tryLater := templates.Notice(
		templates.LoginNoticeID,
		templates.NoticeError,
		tr("error"),
		tr("try_later"),
	)

	if err := database.UpsertUser(h.DB(), ctx, sessionUser); err != nil {
		h.Log().Debug("error upserting user", "error", err)
		tryLater.Render(ctx, w)
		return
	}

	if err := h.Sessions().Set(
		ctx, w,
		sessionUser,
		sessionData,
	); err != nil {
		h.Log().Debug("error setting session", "error", err)
		tryLater.Render(ctx, w)
		return
	}

	h.Log().Debug(
		"logged user in",
		"userID", sessionUser,
		"sessionData", sessionData,
	)

	templates.Redirect(
		config.Endpoints[config.ProtectedPath],
	).Render(ctx, w)
}
