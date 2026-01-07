package router

import (
	"embed"
	"net/http"

	"app/config"
	"app/config/routes"
	"app/handler"
	"app/middleware"
)

func publicEndpoints(h *handler.Handler) []endpoint {
	return []endpoint{
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.Root],
			handler: http.HandlerFunc(h.HomePage),
		},
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.Login],
			handler: http.HandlerFunc(h.LoginPage),
		},
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.About],
			handler: http.HandlerFunc(h.AboutPage),
		},
	}
}

func assetsEndpoints(fs embed.FS) []endpoint {
	root := http.FS(fs)

	handler := http.StripPrefix(
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
			handler: cache(handler),
		},
	}
}
