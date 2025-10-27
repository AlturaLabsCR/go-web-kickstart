// Package handlers implements rendering functions for endpoints
package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"app/i18n"
	"app/sessions"
	"app/utils/smtp"
)

type User struct {
	Email string
}

type Handler struct {
	params     HandlerParams
	Translator func(*http.Request) func(string) string
	Sessions   *sessions.Store[User]
}

type HandlerParams struct {
	Production   bool
	Logger       *slog.Logger
	Database     *sql.DB
	Locales      map[string]map[string]string
	SMTPAuth     smtp.AuthParams
	ServerSecret string
}

func New(params HandlerParams) *Handler {
	sessions := sessions.New[User](sessions.StoreParams{
		CookieName:     "session",
		CookiePath:     "/",
		CookieSameSite: http.SameSiteStrictMode,
		CookieTTL:      24 * time.Hour,
		JWTSecret:      params.ServerSecret,
	})

	translator := i18n.New(params.Locales).TranslateHTTPRequest

	return &Handler{
		params:     params,
		Translator: translator,
		Sessions:   sessions,
	}
}

func (h *Handler) Prod() bool {
	return h.params.Production
}

func (h *Handler) DB() *sql.DB {
	return h.params.Database
}

func (h *Handler) Log() *slog.Logger {
	return h.params.Logger
}

func (h *Handler) SMTPClient() *smtp.Auth {
	return smtp.Client(h.params.SMTPAuth)
}
