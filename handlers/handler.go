// Package handlers implements rendering functions for endpoints
package handlers

import (
	"log/slog"
	"net/http"

	"app/database"
	"app/i18n"
)

type Handler struct {
	params *HandlerParams
}

type HandlerParams struct {
	Production     bool
	Logger         *slog.Logger
	Database       database.Database
	TranslatorFunc i18n.HTTPTranslatorFunc
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

func (h *Handler) Translator(r *http.Request) func(string) string {
	return h.params.TranslatorFunc(r)
}
