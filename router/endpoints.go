package router

import (
	"embed"
	"net/http"

	"app/config"
	"app/config/routes"
	"app/handler"
	"app/middleware"
)

func protectedEndpoints(h *handler.Handler) []endpoint {
	return []endpoint{
		{http.MethodGet, routes.Map[routes.Logout], wrap(h.Validate, h.Logout)},
		{http.MethodGet, routes.Map[routes.Protected], wrap(h.Validate, h.ProtectedPage)},
		{http.MethodGet, routes.Map[routes.ProtectedUser], wrap(h.Validate, h.ProtectedPage)},
		{http.MethodPost, routes.Map[routes.ProtectedUser], wrap(h.Validate, h.ProtectedUpdateUser)},
		{http.MethodGet, routes.Map[routes.ProtectedAdmin], wrap(h.Validate, h.ProtectedAdmin)},
		{http.MethodDelete, routes.Map[routes.ProtectedUser] + "{id}", wrap(h.Validate, h.ProtectedDeleteUser)},
	}
}

func publicEndpoints(h *handler.Handler) []endpoint {
	return []endpoint{
		{http.MethodGet, routes.Map[routes.Login], h.LoginPage},
		{http.MethodGet, routes.Map[routes.About], h.AboutPage},
		{http.MethodPost, routes.Map[routes.GoogleAuth], h.LoginWithGoogle},
		{http.MethodGet, routes.Map[routes.FacebookAuth], h.LoginWithFacebook},
		{
			http.MethodGet,
			routes.Map[routes.Root],
			func(w http.ResponseWriter, r *http.Request) {
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

func wrap(m middleware.Middleware, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m(http.HandlerFunc(h)).ServeHTTP(w, r)
	}
}
