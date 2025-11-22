// Package router implements routing logic to the corresponding handlers
package router

import (
	"net/http"

	"app/handlers"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Router struct {
	routes []Route
}

func Init(h *handlers.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /", h.Home)

	return router
}
