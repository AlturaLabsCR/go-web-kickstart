package config

type Endpoint int

const (
	RootPath = iota
	AssetsPath
	HomePath
	RegisterPath
	LoginPath
)

var Endpoints = map[Endpoint]string{
	RootPath:     "/",
	AssetsPath:   "/assets/",
	HomePath:     "/home",
	RegisterPath: "/register",
	LoginPath:    "/login",
}

func initEndpoints() {
	if prefix := Environment[EnvRootPrefix]; prefix != "" {
		for key := range Endpoints {
			Endpoints[key] = prefix + Endpoints[key]
		}
	}
}
