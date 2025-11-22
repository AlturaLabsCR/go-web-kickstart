// Package handlers implements rendering functions for endpoints
package handlers

import (
	"log/slog"
	"net/http"

	"app/i18n"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	params *HandlerParams
}

type HandlerParams struct {
	Production     bool
	Logger         *slog.Logger
	Database       *pgxpool.Pool
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

func (h *Handler) DB() *pgxpool.Pool {
	return h.params.Database
}

func (h *Handler) Translator(r *http.Request) func(string) string {
	return h.params.TranslatorFunc(r)
}
