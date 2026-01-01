package i18n

import (
	"net/http"
)

type HTTPTranslatorFunc func(*http.Request) func(string) string

func (t *Translator) RequestTranslator(r *http.Request) func(string) string {
	reqLangs := RequestLanguages(r)
	best := t.preferredLanguage(reqLangs)

	return func(key string) string {
		return t.Translate(best, key)
	}
}

func RequestLanguages(r *http.Request) []Language {
	if c, err := r.Cookie(langCookieKey); err == nil {
		return ParseLanguages(c.Value)
	}

	header := r.Header.Get(langHeaderKey)
	return ParseLanguages(header)
}
