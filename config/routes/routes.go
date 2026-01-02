// Package routes
package routes

type Route int

const (
	Root Route = iota
	Assets
	Home
)

var Map = map[Route]string{
	Root:   "/",
	Assets: "/assets/",
	Home:   "/home",
}

func PrefixEndpoints(prefix string) {
	if prefix != "" {
		for key := range Map {
			Map[key] = prefix + Map[key]
		}
	}
}
