// Package router
package router

import (
	"embed"
	"net/http"

	"app/handler"
	"app/middleware"
)

func Init(h *handler.Handler, fs embed.FS) http.Handler {
	mux := http.NewServeMux()

	registerRoutes(mux, h, fs)

	globalMiddleware := middleware.Stack(
		middleware.Gzip,
		h.LogRequest,
	)

	return globalMiddleware(mux)
}
