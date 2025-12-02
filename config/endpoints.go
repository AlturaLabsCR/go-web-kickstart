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
	AuthWithFacebookPath
)

var Endpoints = map[Endpoint]string{
	RootPath:             "/",
	AssetsPath:           "/assets/",
	HomePath:             "/home",
	LoginPath:            "/login",
	LogoutPath:           "/logout",
	ProtectedPath:        "/protected",
	AuthWithGooglePath:   "/auth/google",
	AuthWithFacebookPath: "/auth/facebook",
}

func initEndpoints() {
	if prefix := Config.App.RootPrefix; prefix != "" {
		for key := range Endpoints {
			Endpoints[key] = prefix + Endpoints[key]
		}
	}
}
