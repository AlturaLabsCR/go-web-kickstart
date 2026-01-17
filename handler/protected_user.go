package handler

import (
	"encoding/json"
	"net/http"

	"app/config"
	"app/templates/protected"
)

func (h *Handler) ProtectedUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, ok := r.Context().Value(SessionData).(*config.SessionData)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
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
		h.Log().Debug("invalid username", "body", "r.Body")
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}

	if err := h.DB().Querier().UpsertUserName(ctx,
		payload.UserName,
		session.UserID,
	); err != nil {
		h.Log().Error("error setting username")
		http.Error(w, "error setting username", http.StatusInternalServerError)
		return
	}

	if err := h.Sess().Revoke(w, r); err != nil {
		h.Log().Error("error revoking old session")
		http.Error(w, "error revoking old session", http.StatusInternalServerError)
		return
	}

	if err := h.Sess().Set(ctx, w, session.UserID, session); err != nil {
		h.Log().Error("error setting new session data")
		http.Error(w, "error setting new session data", http.StatusInternalServerError)
		return
	}

	tr := h.Tr(r)

	userMeta, err := h.DB().Querier().GetUserMeta(ctx, session.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Error(w, "error getting user meta", http.StatusInternalServerError)
		return
	}

	if err := protected.ProtectedUser(tr, userMeta, session).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}
