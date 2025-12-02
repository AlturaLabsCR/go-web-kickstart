// Package router implements routing logic to the corresponding handlers
package router

import (
	"embed"
	"net/http"

	"app/config"
	"app/handlers"
	"app/middleware"
)

type endpoint struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func Init(h *handlers.Handler, static embed.FS) http.Handler {
	mux := http.NewServeMux()

	loadRoutes(mux, loadEndpoints(h, static))

	stack := middleware.Stack(
		h.LogRequest,
	)

	return stack(mux)
}

func loadEndpoints(h *handlers.Handler, static embed.FS) []endpoint {
	return []endpoint{
		{
			method: http.MethodGet,
			path:   config.Endpoints[config.AssetsPath],
			handler: h.CachePolicy(
				handlers.MaybeGzip(
					http.StripPrefix(
						config.Config.App.RootPrefix,
						http.FileServer(http.FS(static)),
					),
				),
			),
		},
		{
			method:  http.MethodGet,
			path:    config.Endpoints[config.RootPath],
			handler: h.HomePage,
		},
		{
			method:  http.MethodGet,
			path:    config.Endpoints[config.LoginPath],
			handler: h.LoginPage,
		},
		{
			method:  http.MethodGet,
			path:    config.Endpoints[config.LogoutPath],
			handler: h.Validate(h.Logout),
		},
		{
			method:  http.MethodGet,
			path:    config.Endpoints[config.ProtectedPath],
			handler: h.Validate(h.ProtectedPage),
		},
		{
			method:  http.MethodPost,
			path:    config.Endpoints[config.AuthWithGooglePath],
			handler: h.LoginUserGoogle,
		},
		{
			method:  http.MethodPost,
			path:    config.Endpoints[config.AuthWithFacebookPath],
			handler: h.LoginUserFacebook,
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
		mux.HandleFunc(pattern, e.handler)
	}
}
