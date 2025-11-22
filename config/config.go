// Package config implements initialization logic for required app parameters
package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	AppTitle = "MyApp"
)

const (
	envPrefix     = "APP_"
	EnvProd       = envPrefix + "PROD"
	EnvPort       = envPrefix + "PORT"
	EnvLog        = envPrefix + "LOG_LEVEL"
	EnvConnstr    = envPrefix + "DB_CONNSTR"
	EnvRootPrefix = envPrefix + "ROOT_PREFIX"
)

var Environment = map[string]string{
	EnvProd:       "0",
	EnvPort:       "8080",
	EnvLog:        "0",
	EnvConnstr:    "",
	EnvRootPrefix: "/",
}

func Init() {
	godotenv.Load()

	for key := range Environment {
		if v := os.Getenv(key); v != "" {
			Environment[key] = v
		}
	}
}
