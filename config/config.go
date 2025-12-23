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
	AppTitle                      = "MyApp"
	defaultPort                   = "8080"
	defaultConnStr                = "data/db.db"
	defaultLogLevel               = "0"
	defaultStorageType            = "local"
	defaultStorageLocalRoot       = "data/storage"
	defaultMaxObjectSize    int64 = 0.1e9 // 100MB
	defaultMaxBucketSize    int64 = 1e9   // 1GB
)

const (
	envPrefix = "APP_"
)

type Configuration struct {
	App           AppConfig
	Database      AppDatabase
	AuthProviders AppAuthProviders
	Storage       AppStorage
	Sessions      AppSessions
}

type AppConfig struct {
	Port       string `env:"PORT"`
	LogLevel   string `env:"LOG_LEVEL"`
	RootPrefix string `env:"ROOT_PREFIX"`
}

type AppDatabase struct {
	ConnString string `env:"DB_CONNSTR"`
}

type AppAuthProviders struct {
	Google   GoogleAuthProvider
	Facebook FacebookAuthProvider
}

type GoogleAuthProvider struct {
	ClientID string `env:"GOOGLE_CLIENT_ID"`
}

type FacebookAuthProvider struct {
	AppID     string `env:"FACEBOOK_APP_ID"`
	AppSecret string `env:"FACEBOOK_APP_SECRET"`
}

type AppStorage struct {
	Type          string `env:"STORAGE_TYPE"`
	Remote        AWS
	Local         Local
	MaxObjectSize int64 `env:"STORAGE_MAX_OBJECT_SIZE"`
	MaxBucketSize int64 `env:"STORAGE_MAX_BUCKET_SIZE"`
}

type AWS struct {
	Region            string `env:"AWS_REGION"`
	Bucket            string `env:"AWS_BUCKET"`
	EndpointURL       string `env:"AWS_ENDPOINT_URL"`
	PublicEndpointURL string `env:"AWS_PUBLIC_ENDPOINT_URL"`
	AccessKeyID       string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey   string `env:"AWS_SECRET_ACCESS_KEY"`
}

type Local struct {
	Root              string `env:"STORAGE_LOCAL_ROOT"`
	PublicEndpointURL string `env:"STORAGE_LOCAL_PUBLIC_ENDPOINT_URL"`
}

type AppSessions struct {
	Secret string `env:"SESSIONS_SECRET"`
}

var Config = Configuration{
	App: AppConfig{
		Port:     defaultPort,
		LogLevel: defaultLogLevel,
	},
	Database: AppDatabase{
		ConnString: defaultConnStr,
	},
	Storage: AppStorage{
		Type: defaultStorageType,
		Local: Local{
			Root: defaultStorageLocalRoot,
		},
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

	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("error loading env files: %v", err))
	}

	overrideWithEnv(envPrefix, &Config)

	overrideAWS(&Config.Storage.Remote)

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

func overrideAWS(target *AWS) {
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
				if err := os.Setenv(tag, val); err != nil {
					panic(fmt.Sprintf("error overriding AWS config: %v", err))
				}
			}
		}
	}
}
