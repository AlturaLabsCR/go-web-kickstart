package router

import (
	"embed"
	"net/http"

	"app/config"
	"app/config/routes"
	"app/handler"
	"app/middleware"
)

func pageEndpoints(h *handler.Handler) []endpoint {
	return []endpoint{
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.Root],
			handler: http.HandlerFunc(h.HomePage),
		},
	}
}

func assetEndpoints(static embed.FS) []endpoint {
	root := http.FS(static)

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
