package i18n

import (
	"slices"
	"strings"
)

// BundleExport is an exportable structure representing a Bundle.
type BundleExport struct {
	ETag             string           `json:"etag" yaml:"etag"`
	FallbackLanguage Tag              `json:"fallback_language" yaml:"fallback_language"`
	Languages        []LanguageExport `json:"languages" yaml:"languages"`
}

// LanguageExport is an exportable structure representing single language translations.
type LanguageExport struct {
	ETag         string       `json:"etag" yaml:"etag"`
	Language     Tag          `json:"language" yaml:"language"`
	Translations Translations `json:"translations" yaml:"translations"`
}

// TranslationsFilterFunc is a function deciding add translation key or not to the exported translations.
type TranslationsFilterFunc func(key string) bool

// FilterByPrefix is a TranslationsFilterFunc keeping only translations with the given key prefix.
func FilterByPrefix(prefix string) TranslationsFilterFunc {
	return func(key string) bool {
		return strings.HasPrefix(key, prefix)
	}
}

// NewLanguageExport returns new LanguageExport for the language from the bundle.
func NewLanguageExport(b *Bundle, language Tag, filters ...TranslationsFilterFunc) LanguageExport {
	if b == nil {
		b = NewEmptyBundle()
	}

	etag := FormatLanguageETag(b.CalcHash(), language)

	translations, ok := b.translations[language]
	if !ok {
		return LanguageExport{
			ETag:         etag,
			Language:     language,
			Translations: make(Translations),
		}
	}

	return LanguageExport{
		ETag:         etag,
		Language:     language,
		Translations: FilterTranslations(translations, filters...),
	}
}

// NewBundleExport creates a new BundleExport instance from a Bundle.
func NewBundleExport(b *Bundle, filters ...TranslationsFilterFunc) BundleExport {
	if b == nil {
		b = NewEmptyBundle()
	}

	container := BundleExport{
		ETag:             b.CalcHash(),
		FallbackLanguage: b.fallbackLanguage,
		Languages:        make([]LanguageExport, 0, len(b.translations)),
	}

	for lang, translations := range b.translations {
		container.Languages = append(container.Languages, LanguageExport{
			Language:     lang,
			Translations: FilterTranslations(translations, filters...),
		})
	}

	slices.SortFunc(container.Languages, func(a, b LanguageExport) int {
		if a.Language == b.Language {
			return 0
		}

		if a.Language.String() < b.Language.String() {
			return -1
		}

		return 1
	})

	return container
}

func FormatLanguageETag(bundleHash string, lang Tag) string {
	return bundleHash + "_" + lang.String()
}

// FilterTranslations filters the Translations using the filtering functions.
func FilterTranslations(translations Translations, filters ...TranslationsFilterFunc) Translations {
	result := make(Translations, len(translations))

	if len(filters) == 0 {
		filters = append(filters, func(key string) bool {
			return true
		})
	}

	for key, translation := range translations {
		for _, filter := range filters {
			if !filter(key) {
				continue
			}

			result[key] = translation
		}
	}

	return result
}
