package config

const (
	AppTitle = "MyApp"

	// App
	defaultPort       = "8080"
	defaultLogLevel   = "0"
	defaultRootPrefix = ""
)

var defaultConfigPaths = []string{
	"/etc/app/config.toml",
	"./config.toml",
}
