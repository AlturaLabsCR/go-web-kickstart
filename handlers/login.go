package handlers

import (
	"encoding/json"
	"net/http"

	"app/templates"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	content := templates.Login(h.Translator(r))
	templates.Base(content).Render(r.Context(), w)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := h.Translator(r)

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		templates.Notice(
			templates.LoginNoticeID,
			templates.NoticeError,
			tr("error"),
			tr("login.bad_email"),
		).Render(ctx, w)
		return
	}
}
