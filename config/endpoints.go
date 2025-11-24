package config

type Endpoint int

const (
	RootPath = iota
	AssetsPath
	HomePath
	LoginPath
	RegisterPath
)

var Endpoints = map[Endpoint]string{
	RootPath:     "/",
	AssetsPath:   "/assets/",
	HomePath:     "/home",
	LoginPath:    "/login",
	RegisterPath: "/register",
}
