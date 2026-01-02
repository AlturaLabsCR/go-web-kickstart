package config

import (
	"os"
	"reflect"
	"strconv"
)

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
