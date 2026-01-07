// Package config
package config

import "time"

const (
	envPrefix = "APP_"
)

var Config = Configuration{
	App: App{
		Port:       defaultPort,
		LogLevel:   defaultLogLevel,
		RootPrefix: defaultRootPrefix,
	},
	Database: Database{
		ConnString: defaultDBConnStr,
	},
	Year: time.Now().Year(),
}

type Configuration struct {
	App      App
	Database Database
	Sessions Sessions
	Year     int
}

type Database struct {
	ConnString string `env:"DB_CONNSTR"`
}

type App struct {
	Port       string `env:"PORT"`
	LogLevel   string `env:"LOG_LEVEL"`
	RootPrefix string `env:"ROOT_PREFIX"`
}

type Sessions struct {
	Secret string `env:"SESSIONS_SECRET"`
}
