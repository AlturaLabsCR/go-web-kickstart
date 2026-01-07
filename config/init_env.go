package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

func InitEnv() {
	for _, name := range defaultConfigPaths {
		if data, err := os.ReadFile(name); err == nil {
			if _, err := toml.Decode(string(data), &Config); err != nil {
				panic(fmt.Sprintf("error decoding config: %v", err))
			}
			break
		}
	}

	_ = godotenv.Load()
	overrideWithEnv(envPrefix, &Config)
	overrideAWS(&Config.Storage.S3)
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

// overrideAWS sets AWS environment variables from the given S3 config.
// For each struct field tagged with `env`, it sets the env var only if:
//   - the env var is not already defined, and
//   - the struct field value is non-empty.
//
// This allows config-file values to act as defaults while preserving
// explicitly provided environment variables.
//
// Used by InitEnv, required by initS3.
func overrideAWS(target *S3) {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envKey := fieldType.Tag.Get("env")
		if envKey == "" {
			continue
		}

		if field.Kind() != reflect.String {
			continue
		}

		if _, exists := os.LookupEnv(envKey); exists {
			continue
		}

		value := field.String()
		if value == "" {
			continue
		}

		_ = os.Setenv(envKey, value)
	}
}
