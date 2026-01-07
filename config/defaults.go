package config

const (
	AppTitle = "MyApp"

	// App
	defaultPort       = "8080"
	defaultLogLevel   = "0"
	defaultRootPrefix = ""

	// DB
	defaultDBConnStr = "./data/db.sqlite"

	// Storage
	defaultStorageFS       = "fs"
	defaultStorageFSRoot   = "./data/storage"
	defaultStorageS3Region = "auto"
	defaultMaxObjectSize   = 1_000_000
	defaultMaxBucketSize   = 1_000_000_000
)

var defaultConfigPaths = []string{
	"/etc/app/config.toml",
	"./config.toml",
}
