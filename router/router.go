// Package router
package router

import (
	"embed"
	"net/http"

	"app/config"
	"app/handler"
	"app/middleware"
)

func Init(h *handler.Handler, fs embed.FS) http.Handler {
	mux := http.NewServeMux()

	registerRoutes(mux, h, fs)

	var globalMiddleware middleware.Middleware

	// do not compress if log level <= 0 (debug mode), because:
	// templ generate --notify-proxy requires the html response to inject
	// the hot-reload script in development.
	if config.Config.App.Env == config.ProductionEnv {
		globalMiddleware = middleware.Stack(middleware.Gzip)
	} else {
		globalMiddleware = middleware.Stack(middleware.Delay, h.LogRequest, middleware.ContentLength)
	}

	return globalMiddleware(mux)
}
