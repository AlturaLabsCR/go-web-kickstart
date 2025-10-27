// Package router implements routing logic to the corresponding handlers
package router

import (
	"net/http"

	"app/handlers"
	"app/middleware"
)

func Routes(h *handlers.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /", h.Home)
	router.HandleFunc("GET /login", h.Login)
	router.HandleFunc("POST /login", h.SendVerification)
	router.HandleFunc("GET /verify", h.Verify)
	router.Handle("GET /{name}", middleware.With(h.Protected, h.RenderName))
	router.HandleFunc("GET /dogs", h.ListDogs)

	return router
}
