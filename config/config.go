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
	Storage: Storage{
		Type: defaultStorageFS,
		FS: FS{
			Root: defaultStorageFSRoot,
		},
		S3: S3{
			Region: defaultStorageS3Region,
		},
		MaxObjectSize: defaultMaxObjectSize,
		MaxBucketSize: defaultMaxBucketSize,
	},
	Year: time.Now().Year(),
}

type Configuration struct {
	App           App
	Database      Database
	Sessions      Sessions
	Storage       Storage
	AuthProviders AuthProviders
	Year          int
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

type Storage struct {
	Type          string `env:"STORAGE_TYPE"`
	S3            S3
	FS            FS
	MaxObjectSize int64 `env:"STORAGE_MAX_OBJECT_SIZE"`
	MaxBucketSize int64 `env:"STORAGE_MAX_BUCKET_SIZE"`
}

type S3 struct {
	Region          string `env:"AWS_REGION"`
	Bucket          string `env:"AWS_BUCKET"`
	EndpointURL     string `env:"AWS_ENDPOINT_URL"`
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
}

type FS struct {
	Root string `env:"STORAGE_LOCAL_ROOT"`
}

type AuthProviders struct {
	Google   Google
	Facebook Facebook
}

type Google struct {
	ClientID string `env:"GOOGLE_CLIENT_ID"`
}

type Facebook struct {
	AppID     string `env:"FACEBOOK_APP_ID"`
	AppSecret string `env:"FACEBOOK_APP_SECRET"`
}
