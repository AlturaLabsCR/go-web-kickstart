// Package handlers implements rendering functions for endpoints
package handlers

import (
	"log/slog"
	"net/http"

	"app/config"
	"app/database"
	"app/i18n"
	"app/sessions"
	"app/storage/s3"
)

type Handler struct {
	params *HandlerParams
}

type HandlerParams struct {
	Production     bool
	Logger         *slog.Logger
	Database       database.Database
	Storage        s3.Storage
	TranslatorFunc i18n.HTTPTranslatorFunc
	Sessions       *sessions.Store[config.SessionData]
	Secret         string
}

func New(params *HandlerParams) *Handler {
	return &Handler{params}
}

func (h *Handler) Production() bool {
	return h.params.Production
}

func (h *Handler) Log() *slog.Logger {
	return h.params.Logger
}

func (h *Handler) DB() database.Database {
	return h.params.Database
}

func (h *Handler) S3() s3.Storage {
	return h.params.Storage
}

func (h *Handler) Translator(r *http.Request) func(string) string {
	return h.params.TranslatorFunc(r)
}

func (h *Handler) Sessions() *sessions.Store[config.SessionData] {
	return h.params.Sessions
}
