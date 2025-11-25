package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"app/config"
	"app/database"
	"app/templates"
)

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	content := templates.Register(h.Translator(r))
	templates.Base(content).Render(r.Context(), w)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := h.Translator(r)

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		templates.Notice(
			templates.RegisterNoticeID,
			templates.NoticeError,
			tr("error"),
			tr("register.bad_email"),
		).Render(ctx, w)
		return
	}

	if _, err := database.InsertUser(
		h.DB(),
		r.Context(),
		req.Email,
	); err != nil {
		h.Log().Error("error registering user", "error", err)

		if errors.Is(err, database.ErrDuplicateEmail) {
			templates.Notice(
				templates.RegisterNoticeID,
				templates.NoticeWarn,
				tr("warn"),
				tr("register.email_exists"),
			).Render(ctx, w)
		} else {
			templates.Notice(
				templates.RegisterNoticeID,
				templates.NoticeError,
				tr("error"),
				tr("register.bad_email"),
			).Render(ctx, w)
		}
		return
	}

	templates.Redirect(config.Endpoints[config.LoginPath]).Render(ctx, w)
}
