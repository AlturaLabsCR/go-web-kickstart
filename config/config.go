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
	envPrefix = "APP_"

	// required
	EnvDriver  = envPrefix + "DB_DRIVER"
	EnvConnStr = envPrefix + "DB_CONNSTR"

	// optional
	EnvProd       = envPrefix + "PROD"
	EnvPort       = envPrefix + "PORT"
	EnvLog        = envPrefix + "LOG_LEVEL"
	EnvRootPrefix = envPrefix + "ROOT_PREFIX"
)

var Environment = map[string]string{
	EnvConnStr: "",

	EnvProd:       "0",
	EnvPort:       "8080",
	EnvLog:        "0",
	EnvRootPrefix: "",
}

func Init() {
	godotenv.Load()

	for key := range Environment {
		if v := os.Getenv(key); v != "" {
			Environment[key] = v
		}
	}
}
