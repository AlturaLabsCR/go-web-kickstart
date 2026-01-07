// Package handler
package handler

import (
	"log/slog"
	"net/http"

	"app/config"
	"app/database"
	"app/i18n"
	"app/sessions"
	"app/storage"
)

type Handler struct {
	params *HandlerParams
}

type HandlerParams struct {
	Logger     *slog.Logger
	Database   database.Database
	Storage    storage.ObjectStorage
	Translator i18n.HTTPTranslatorFunc
	Sessions   *sessions.Store[config.SessionData]
}

func New(params *HandlerParams) (*Handler, error) {
	if params.Logger == nil {
		return nil, errStr("logger is nil")
	}
	if params.Database == nil {
		return nil, errStr("database interface is nil")
	}
	if params.Storage == nil {
		return nil, errStr("storage interface is nil")
	}
	if params.Translator == nil {
		return nil, errStr("translator func is nil")
	}
	if params.Sessions == nil {
		return nil, errStr("sessions interface is nil")
	}

	return &Handler{params}, nil
}

func (h *Handler) Log() *slog.Logger {
	return h.params.Logger
}

func (h *Handler) DB() database.Database {
	return h.params.Database
}

func (h *Handler) S3() storage.ObjectStorage {
	return h.params.Storage
}

func (h *Handler) Tr(r *http.Request) func(string) string {
	return h.params.Translator(r)
}

func (h *Handler) Sess() *sessions.Store[config.SessionData] {
	return h.params.Sessions
}
