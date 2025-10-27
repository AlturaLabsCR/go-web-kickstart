package config

import "app/i18n"

func InitLocales() map[string]map[string]string {
	return map[string]map[string]string{
		"es": i18n.ES,
		"en": i18n.EN,
	}
}
