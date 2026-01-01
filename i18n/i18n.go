// Package i18n provides simple locale selection and translation utilities
package i18n

import (
	"sort"
	"strings"
)

type Language struct {
	Tag     string
	Variant string
	Weight  float64
}

type Locale map[string]string

type Locales map[string]Locale

type Translator struct {
	locales  Locales
	fallback string
}

func New(locales Locales, fallback string) (*Translator, error) {
	if len(locales) < 1 {
		return nil, ErrEmptyLocales
	}

	if _, ok := locales[fallback]; !ok {
		return nil, ErrInvalidFallback
	}

	return &Translator{locales, fallback}, nil
}

func (t *Translator) Translate(lang, key string) string {
	if loc, ok := t.locales[lang]; ok {
		if val, ok := loc[key]; ok {
			return val
		}
	}
	return key
}

func ParseLanguages(s string) []Language {
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return []Language{}
	}

	langs := []Language{}

	for langStr := range strings.SplitSeq(s, ",") {
		if l := parseLanguage(langStr); l.Tag != "" {
			langs = append(langs, l)
		}
	}

	sort.SliceStable(langs, func(i, j int) bool {
		return langs[i].Weight > langs[j].Weight
	})

	return langs
}
