package router

import (
	"embed"
	"net/http"

	"app/handler"
)

type endpoint struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func registerRoutes(
	mux *http.ServeMux,
	h *handler.Handler,
	fs embed.FS,
) {
	for _, e := range endpoints(h, fs) {
		pattern := e.path
		if e.method != "" {
			pattern = e.method + " " + e.path
		}
		mux.HandleFunc(pattern, e.handler)
	}
}

func endpoints(h *handler.Handler, fs embed.FS) []endpoint {
	var endpoints []endpoint

	endpoints = append(endpoints, publicEndpoints(h)...)
	endpoints = append(endpoints, assetsEndpoints(fs)...)

	return endpoints
}
