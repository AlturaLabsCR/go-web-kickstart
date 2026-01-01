package i18n

import (
	"slices"
	"sort"
	"strconv"
	"strings"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const (
	ErrEmptyLocales    = errStr("empty locales")
	ErrInvalidFallback = errStr("fallback key does not exist")

	langHeaderKey = "Accept-Language"
	langCookieKey = "lang"
)

func parseLanguage(s string) Language {
	base, qpart, hasQ := strings.Cut(s, ";")

	if base == "" {
		return Language{}
	}

	tag, variant, _ := strings.Cut(base, "-")

	weight := 1.0

	if hasQ {
		if _, val, ok := strings.Cut(qpart, "q="); ok {
			if w, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err == nil {
				weight = w
			} else {
				weight = 0
			}
		}
	}

	return Language{
		Tag:     tag,
		Variant: variant,
		Weight:  weight,
	}
}

func (t *Translator) preferredLanguage(req []Language) string {
	if len(t.locales) == 0 {
		return ""
	}

	supported := make([]string, 0, len(t.locales))
	for lang := range t.locales {
		supported = append(supported, strings.ToLower(string(lang)))
	}
	sort.Strings(supported)

	for _, r := range req {
		full := r.Tag
		if r.Variant != "" {
			full = full + "-" + r.Variant
		}

		// exact match: "es-419"
		if slices.Contains(supported, full) {
			return full
		}

		// tag-only match: "es"
		if slices.Contains(supported, r.Tag) {
			return r.Tag
		}
	}

	return t.fallback
}
