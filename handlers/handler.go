// Package handlers implements rendering functions for endpoints
package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"app/utils/smtp"
)

type Handler struct {
	params     HandlerParams
	Translator func(*http.Request) func(string) string
	SMTPClient *smtp.Auth
}

type HandlerParams struct {
	Production bool
	DB         *sql.DB
	Logger     *slog.Logger
}

func New(params HandlerParams, translatorFunc func(*http.Request) func(string) string, smtpParams smtp.AuthParams) *Handler {
	return &Handler{
		params:     params,
		Translator: translatorFunc,
		SMTPClient: smtp.Client(smtpParams),
	}
}

func (h *Handler) Prod() bool {
	return h.params.Production
}

func (h *Handler) DB() *sql.DB {
	return h.params.DB
}

func (h *Handler) Log() *slog.Logger {
	return h.params.Logger
}
