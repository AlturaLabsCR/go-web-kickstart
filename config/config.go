// Package config implements initialization logic for required app parameters
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

const (
	AppTitle                   = "MyApp"
	defaultPort                = "8080"
	defaultConnStr             = "data/db.db"
	defaultLogLevel            = "0"
	defaultStorageType         = "local"
	defaultStoragePath         = "data/storage"
	defaultMaxObjectSize int64 = 0.1e9 // 100MB
	defaultMaxBucketSize int64 = 1e9   // 1GB
)

const (
	envPrefix = "APP_"
)

type Configuration struct {
	App           AppConfig
	DB            AppDatabase
	AuthProviders AppAuthProviders
	Storage       AppStorage
}

type AppConfig struct {
	Port       string `env:"PORT"`
	LogLevel   string `env:"LOG_LEVEL"`
	RootPrefix string `env:"ROOT_PREFIX"`
	Secret     string `env:"SECRET"`
}

type AppDatabase struct {
	ConnString string `env:"DB_CONNSTR"`
}

type AppAuthProviders struct {
	Google   GoogleCredentials
	Facebook FacebookCredentials
}

type GoogleCredentials struct {
	ClientID string `env:"GOOGLE_CLIENT_ID"`
}

type FacebookCredentials struct {
	AppID     string `env:"FACEBOOK_APP_ID"`
	AppSecret string `env:"FACEBOOK_APP_SECRET"`
}

type AppStorage struct {
	Type          string `env:"STORAGE_TYPE"`
	LocalRoot     string `env:"STORAGE_ROOT"`
	RemoteBucket  string `env:"STORAGE_BUCKET"`
	MaxObjectSize int64  `env:"STORAGE_MAX_OBJECT_SIZE"`
	MaxBucketSize int64  `env:"STORAGE_MAX_BUCKET_SIZE"`
	AWS           AWSCredentials
}

type AWSCredentials struct {
	AwsRegion          string `env:"AWS_REGION"`
	AwsEndpointURL     string `env:"AWS_ENDPOINT_URL"`
	AwsAccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
}

var Config = Configuration{
	App: AppConfig{
		Port:     defaultPort,
		LogLevel: defaultLogLevel,
	},
	DB: AppDatabase{
		ConnString: defaultConnStr,
	},
	Storage: AppStorage{
		Type:          defaultStorageType,
		LocalRoot:     defaultStoragePath,
		MaxObjectSize: defaultMaxObjectSize,
		MaxBucketSize: defaultMaxBucketSize,
	},
}

var configPaths = []string{
	"/etc/app/config.toml",
	"./config.toml",
}

func Init() {
	for _, conf := range configPaths {
		if data, err := os.ReadFile(conf); err == nil {
			if _, err := toml.Decode(string(data), &Config); err != nil {
				panic(fmt.Sprintf("error decoding config: %v", err))
			}
			break
		}
	}

	godotenv.Load()

	overrideWithEnv(envPrefix, &Config)

	overrideAWS(&Config.Storage.AWS)

	initEndpoints()
}

func overrideWithEnv(prefix string, target any) {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			overrideWithEnv(prefix, field.Addr().Interface())
			continue
		}

		tag, ok := fieldType.Tag.Lookup("env")
		if !ok {
			continue
		}

		envKey := prefix + tag
		envVal, exists := os.LookupEnv(envKey)
		if !exists || envVal == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envVal)
		case reflect.Int64:
			val, err := strconv.ParseInt(envVal, 10, 64)
			if err != nil {
				panic("failed to parse int")
			}
			field.SetInt(val)
		}
	}
}

func overrideAWS(target *AWSCredentials) {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag, ok := fieldType.Tag.Lookup("env")
		if !ok {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if val := field.String(); val != "" {
				os.Setenv(tag, val)
			}
		}
	}
}
