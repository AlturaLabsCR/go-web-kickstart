package config

type Endpoint int

const (
	AssetsPath = iota
	HomePath
	LoginPath
)

var Endpoints = map[Endpoint]string{
	AssetsPath: "assets/",
	HomePath:   "home",
	LoginPath:  "login",
}
