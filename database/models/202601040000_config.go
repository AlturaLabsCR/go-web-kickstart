package models

const (
	ConfigCacheScopePrefix = "config:"

	ConfigInitialized     = "config.initialized"
	ConfigInitializedTrue = "true"
)

type Config struct {
	Name  string
	Value string
}
