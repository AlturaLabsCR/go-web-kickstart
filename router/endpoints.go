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
			path:    routes.Map[routes.Login],
			handler: h.LoginPage,
		},
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.About],
			handler: h.AboutPage,
		},
		{
			method:  http.MethodPost,
			path:    routes.Map[routes.GoogleAuth],
			handler: h.LoginWithGoogle,
		},
		{
			method:  http.MethodGet,
			path:    routes.Map[routes.FacebookAuth],
			handler: h.AboutPage,
		},
		{
			method: http.MethodGet,
			path:   routes.Map[routes.Root],
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != routes.Map[routes.Root] {
					http.Redirect(w, r, routes.Map[routes.Root], http.StatusFound)
					return
				}
				h.HomePage(w, r)
			},
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
			handler: cache(handler).(http.HandlerFunc),
		},
	}
}
