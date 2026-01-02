package config

import "app/config/routes"

func InitRoutes() {
	if prefix := Config.App.RootPrefix; prefix != "" {
		routes.PrefixEndpoints(prefix)
	}
}
