package config

import (
	"fmt"

	"app/i18n"
)

func InitTranslator() i18n.HTTPTranslatorFunc {
	locales := i18n.Locales{
		i18n.ESKey: i18n.ES,
		i18n.ENKey: i18n.EN,
	}

	translator, err := i18n.New(locales, i18n.ENKey)
	if err != nil {
		panic(fmt.Sprintf("error initializing translator: %v", err))
	}

	return translator.RequestTranslator
}
