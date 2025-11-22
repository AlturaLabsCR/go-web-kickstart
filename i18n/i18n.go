// Package i18n provides simple locale selection and translation utilities.
package i18n

import (
	"net/http"
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
	langHeaderKey = "Accept-Language"
	langCookieKey = "lang"

	ErrEmptyLocales    = errStr("passed locales map is empty")
	ErrInvalidFallback = errStr("passed fallback key does not exist in locales map")
)

// HTTPTranslatorFunc represents a function returning a function that translates
// a key, given a request
type HTTPTranslatorFunc func(*http.Request) func(string) string

// Language represents a parsed Accept-Language entry.
type Language struct {
	Tag     string  // e.g. "en"
	Variant string  // e.g. "US"
	Weight  float64 // q-value, 1.0 = highest
}

// Locale is a set of translation key-value pairs.
type Locale map[string]string

// Locales is a set of locale key-value pairs.
type Locales map[string]Locale

// Translator loads locales and performs language lookup.
type Translator struct {
	locales  Locales
	fallback string
}

// New creates a new Translator with the given locales.
// Example keys: "en", "es", "es-419"
func New(locales Locales, fallback string) (*Translator, error) {
	if len(locales) < 1 {
		return nil, ErrEmptyLocales
	}

	if _, ok := locales[fallback]; !ok {
		return nil, ErrInvalidFallback
	}

	return &Translator{locales, fallback}, nil
}

// RequestTranslator returns a function that translates keys according to the
// best language available for the incoming http.Request.
func (t *Translator) RequestTranslator(r *http.Request) func(string) string {
	reqLangs := requestLanguages(r)
	best := t.preferredLanguage(reqLangs)

	return func(key string) string {
		return t.Translate(best, key)
	}
}

// Translate returns a translated key in the given language. If missing, falls
// back to returning the key itself.
func (t *Translator) Translate(lang, key string) string {
	if loc, ok := t.locales[lang]; ok {
		if val, ok := loc[key]; ok {
			return val
		}
	}
	return key
}

// preferredLanguage selects the best matching language from the translator's
// supported locales.
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

// requestLanguages obtains language preferences from cookies or Accept-Language
// header.
func requestLanguages(r *http.Request) []Language {
	if c, err := r.Cookie(langCookieKey); err == nil {
		return parseLanguages(c.Value)
	}

	header := r.Header.Get(langHeaderKey)
	return parseLanguages(header)
}

// parseLanguages parses an Accept-Language header into a slice of Language
// objects, sorted by descending weight.
func parseLanguages(s string) []Language {
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return nil
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

// parseLanguage parses a single language entry such as "es-419;q=0.9".
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
