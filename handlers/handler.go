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

func (h *Handler) Translator(r *http.Request) func(string) string {
	return h.params.TranslatorFunc(r)
}
