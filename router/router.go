// Package router implements routing logic to the corresponding handlers
package router

import (
	"net/http"

	"app/handlers"
)

type endpoint struct {
	method  string
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

func Init(h *handlers.Handler) *http.ServeMux {
	router := http.NewServeMux()

	endpoints := []endpoint{
		{method: "", path: "", handler: nil},
	}

	loadRoutes(router, endpoints)

	return router
}

func loadRoutes(router *http.ServeMux, endpoints []endpoint) {
	for _, e := range endpoints {
		var pattern string
		if e.method != "" {
			pattern = e.method + " " + e.path
		} else {
			pattern = e.path
		}
		router.HandleFunc(pattern, e.handler)
	}
}
