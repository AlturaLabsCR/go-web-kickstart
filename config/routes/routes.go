// Package routes
package routes

type Route int

const (
	Root Route = iota
	Assets
	Login
	About
)

var Map = map[Route]string{
	Root:   "/",
	Assets: "/assets/",
	Login:  "/login/",
	About:  "/about/",
}

func PrefixEndpoints(prefix string) {
	if prefix != "" {
		for key := range Map {
			Map[key] = prefix + Map[key]
		}
	}
}
