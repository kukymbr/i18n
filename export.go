package i18n

import "slices"

// BundleExport is an exportable structure representing a Bundle.
type BundleExport struct {
	FallbackLanguage Tag              `json:"fallback_language" yaml:"fallback_language"`
	Languages        []LanguageExport `json:"languages" yaml:"languages"`
}

// LanguageExport is an exportable structure representing single language translations.
type LanguageExport struct {
	Language     Tag          `json:"language" yaml:"language"`
	Translations Translations `json:"translations" yaml:"translations"`
}

// NewBundleExport creates a new BundleExport instance from a Bundle.
func NewBundleExport(b *Bundle) BundleExport {
	container := BundleExport{
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
