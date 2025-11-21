package config

import "app/i18n"

func InitLocales() map[string]i18n.Locale {
	return map[string]i18n.Locale{
		"es": i18n.ES,
		"en": i18n.EN,
	}
}
