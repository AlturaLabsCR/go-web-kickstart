// Package router
package router

import (
	"embed"
	"net/http"

	"app/config"
	"app/config/routes"
	"app/handler"
	"app/middleware"
)

type endpoint struct {
	method  string
	path    string
	handler http.Handler
}

func Init(h *handler.Handler, static embed.FS) http.Handler {
	mux := http.NewServeMux()
	loadRoutes(mux, loadEndpoints(h, static))

	globalMiddleware := middleware.Stack(
		middleware.Gzip,
		h.LogRequest,
	)

	return globalMiddleware(mux)
}

func loadEndpoints(h *handler.Handler, fsys embed.FS) []endpoint {
	root := http.FS(fsys)
	static := http.StripPrefix(
		config.Config.App.RootPrefix,
		http.FileServer(root),
	)

	cache := middleware.Stack(
		middleware.Cache(
			middleware.CachePolicy{
				Enabled: true,
				Public:  true,
			},
		),
	)

	return []endpoint{
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.Assets],
			handler: cache(static),
		},
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.Root],
			handler: http.HandlerFunc(h.HomePage),
		},
	}
}

func loadRoutes(mux *http.ServeMux, endpoints []endpoint) {
	for _, e := range endpoints {
		var pattern string
		if e.method != "" {
			pattern = e.method + " " + e.path
		} else {
			pattern = e.path
		}
		mux.Handle(pattern, e.handler)
	}
}
