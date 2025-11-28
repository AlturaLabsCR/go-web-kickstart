package config

type Endpoint int

const (
	RootPath = iota
	AssetsPath
	HomePath
	LoginPath
	LogoutPath
	ProtectedPath
	AuthWithGooglePath
)

var Endpoints = map[Endpoint]string{
	RootPath:           "/",
	AssetsPath:         "/assets/",
	HomePath:           "/home",
	LoginPath:          "/login",
	LogoutPath:         "/logout",
	ProtectedPath:      "/protected",
	AuthWithGooglePath: "/auth/google",
}

func initEndpoints() {
	if prefix := Environment[EnvRootPrefix]; prefix != "" {
		for key := range Endpoints {
			Endpoints[key] = prefix + Endpoints[key]
		}
	}
}
