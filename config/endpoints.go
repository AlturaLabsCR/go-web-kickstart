package config

type Endpoint int

const (
	RootPath = iota
	AssetsPath
	HomePath
	LoginPath
	ProtectedPath
)

var Endpoints = map[Endpoint]string{
	RootPath:      "/",
	AssetsPath:    "/assets/",
	HomePath:      "/home",
	LoginPath:     "/login",
	ProtectedPath: "/protected",
}

func initEndpoints() {
	if prefix := Environment[EnvRootPrefix]; prefix != "" {
		for key := range Endpoints {
			Endpoints[key] = prefix + Endpoints[key]
		}
	}
}
