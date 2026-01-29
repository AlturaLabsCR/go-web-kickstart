package handler

import (
	"encoding/json"
	"net/http"

	"app/database/models"
	"app/i18n"
	"app/templates/protected"
)

func (h *Handler) ProtectedUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionData, ok := h.Sess().Data(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	sessionAttrs, ok := h.Sess().Attrs(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	user, err := h.DB().Querier().GetUserMeta(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Error(w, "error getting user meta", http.StatusBadRequest)
		return
	}

	if !models.HasPermission(user.Perms, "perm.change_name") {
		h.Log().Error("doesnt have permission to change name")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload struct {
		UserName string `json:"onboarding-input-name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.Log().Debug("invalid json", "body", "r.Body")
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.UserName == "" {
		h.Log().Debug("invalid username", "body", r.Body)
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}

	if err := h.DB().Querier().UpsertUserName(ctx,
		payload.UserName,
		sessionData.UserID,
	); err != nil {
		h.Log().Error("error setting username")
		http.Error(w, "error setting username", http.StatusInternalServerError)
		return
	}

	tr := h.Tr(r)

	userMeta, err := h.DB().Querier().GetUserMeta(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Error(w, "error getting user meta", http.StatusInternalServerError)
		return
	}

	userSessions, err := h.Sess().AttrsByUser(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting sessions attrs", "error", err)
		http.Error(w, "error getting session attrs", http.StatusInternalServerError)
		return
	}

	locale := "es"

	langs := i18n.RequestLanguages(r)
	if len(langs) > 0 && langs[0].Tag != "" {
		locale = langs[0].Tag
	}

	params := protected.ProtectedParams{
		User:     userMeta,
		Attrs:    sessionAttrs,
		Data:     sessionData,
		Sessions: userSessions,
		Active:   r.URL.Path,
		Locale:   locale,
	}

	if err := protected.ProtectedUser(tr, params).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
