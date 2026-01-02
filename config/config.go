// Package config
package config

import (
	"fmt"
	"os"

	"app/config/routes"

	"github.com/BurntSushi/toml"
)

var Config = Configuration{
	App: App{
		Port:       defaultPort,
		LogLevel:   defaultLogLevel,
		RootPrefix: defaultRootPrefix,
	},
}

func Init() {
	for _, name := range defaultConfigPaths {
		if data, err := os.ReadFile(name); err == nil {
			if _, err := toml.Decode(string(data), &Config); err != nil {
				panic(fmt.Sprintf("error decoding config: %v", err))
			}
			break
		}
	}

	overrideWithEnv(envPrefix, &Config)

	if prefix := Config.App.RootPrefix; prefix != "" {
		routes.PrefixEndpoints(prefix)
	}
}
