// Package config implements initialization logic for required app parameters
package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

const (
	AppTitle        = "MyApp"
	defaultPort     = "8080"
	defaultConnStr  = "./db.db"
	defaultLogLevel = "0"
)

const (
	envPrefix = "APP_"
)

type Configuration struct {
	App         AppConfig
	Credentials AppCredentials
}

type AppConfig struct {
	Port       string `env:"PORT"`
	ConnString string `env:"DB_CONNSTR"`
	LogLevel   string `env:"LOG_LEVEL"`
	RootPrefix string `env:"ROOT_PREFIX"`
	Secret     string `env:"SECRET"`
}

type AppCredentials struct {
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

var Config = Configuration{
	App: AppConfig{
		Port:       defaultPort,
		ConnString: defaultConnStr,
		LogLevel:   defaultLogLevel,
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
		}
	}
}
