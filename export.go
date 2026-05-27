package i18n

import "slices"

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

// NewLanguageExport returns new LanguageExport for the language from the bundle.
func NewLanguageExport(b *Bundle, language Tag) LanguageExport {
	if b == nil {
		b = NewEmptyBundle()
	}

	translations, ok := b.translations[language]
	if !ok {
		return LanguageExport{
			ETag:         FormatLanguageETag(b.CalcHash(), language),
			Language:     language,
			Translations: make(Translations),
		}
	}

	return LanguageExport{
		Language:     language,
		Translations: translations,
	}
}

// NewBundleExport creates a new BundleExport instance from a Bundle.
func NewBundleExport(b *Bundle) BundleExport {
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
			Translations: translations,
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
