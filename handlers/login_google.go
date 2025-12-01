package handlers

import (
	"net/http"

	"app/auth"
	"app/config"
	"app/database"

	"github.com/mileusna/useragent"
)

func (h *Handler) LoginUserGoogle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// if the session is valid, redirect
	if _, err := h.Sessions().Validate(w, r); err == nil {
		http.Redirect(w, r, config.Endpoints[config.ProtectedPath], http.StatusSeeOther)
		return
	}

	sessionUser, err := auth.GetGoogleID(r, config.Environment[config.EnvGoogleClientID])
	if err != nil {
		h.Log().Debug("error getting sessionUser", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// from now on the user is valid

	ua := useragent.Parse(r.UserAgent())
	sessionData := config.SessionData{
		OS:   ua.OS,
		Name: ua.Name,
	}

	if err := database.UpsertUser(h.DB(), ctx, sessionUser); err != nil {
		h.Log().Error("error upserting user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.Sessions().Set(
		ctx, w,
		sessionUser,
		sessionData,
	); err != nil {
		h.Log().Debug("error setting session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log().Debug(
		"logged user in",
		"userID", sessionUser,
		"sessionData", sessionData,
	)

	http.Redirect(w, r, config.Endpoints[config.ProtectedPath], http.StatusSeeOther)
}
