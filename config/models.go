package config

const (
	envPrefix = "TE_"
)

type Configuration struct {
	App      App
	Sessions Sessions
}

type App struct {
	Port       string `env:"PORT"`
	LogLevel   string `env:"LOG_LEVEL"`
	RootPrefix string `env:"ROOT_PREFIX"`
}

type Sessions struct {
	Secret string `env:"SESSIONS_SECRET"`
}
