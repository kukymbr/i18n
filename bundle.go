package i18n

import (
	"fmt"
	"strings"

	"github.com/kukymbr/i18n/internal/tagsparser"
	"golang.org/x/text/language"
)

// NewBundle creates new Bundle instance.
func NewBundle[T Language](fallbackLanguage T, sources ...BundleSource) (*Bundle, error) {
	b := &Bundle{fallbackLanguage: Lang(fallbackLanguage)}

	b.translations = make(map[language.Tag]Translations)

	for _, source := range sources {
		if err := source(b); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Bundle is an i18n translations bundle.
type Bundle struct {
	fallbackLanguage language.Tag
	translations     map[language.Tag]Translations
}

// Translate finds a translation for a key.
func (b *Bundle) Translate(lang language.Tag, key string, tplData ...any) string {
	var data any
	if len(tplData) > 0 {
		data = tplData[0]
	}

	return b.translate(lang, key, data)
}

// T is a short alias for a Translate.
func (b *Bundle) T(lang language.Tag, key string, tplData ...any) string {
	return b.Translate(lang, key, tplData...)
}

// TranslateStruct updated fields of the given structure with a translated representation.
// The structure argument must be a pointer to a non-nil structure variable.
//
// The `i18n:"field.key"` tag format is expected to get a field's translation key;
// if no `i18n` found, field's value is used as a key.
// Add `i18n:"-"` tag to skip field's translation.
// Only string values are affected.
func (b *Bundle) TranslateStruct(lang language.Tag, structure any, tplData ...any) error {
	err := tagsparser.ParseTags(structure, func(s string) string {
		return b.Translate(lang, s, tplData...)
	})
	if err != nil {
		return fmt.Errorf("translate structure: %w", err)
	}

	return nil
}

// Translate finds a translation for a key.
func (b *Bundle) translate(lang language.Tag, key string, tplData any) string {
	if lang == language.Und {
		lang = b.fallbackLanguage
	}

	keys := []string{
		key,
		strings.ToLower(key),
	}

	for _, k := range keys {
		text, ok := b.getTranslation(lang, k)
		if ok {
			return prepareText(k, text, tplData)
		}
	}

	if lang != b.fallbackLanguage {
		return b.translate(b.fallbackLanguage, key, tplData)
	}

	return prepareText(key, key, tplData)
}

func (b *Bundle) getTranslation(lang language.Tag, key string) (string, bool) {
	text, ok := b.translations[lang][key]

	return text, ok
}

func (b *Bundle) addTranslations(lang language.Tag, translations Translations) {
	for key, text := range translations {
		b.addTranslation(lang, key, text)
	}
}

func (b *Bundle) addTranslation(lang language.Tag, key string, text string) {
	if lang == language.Und {
		lang = b.fallbackLanguage
	}

	if _, ok := b.translations[lang]; !ok {
		b.translations[lang] = make(Translations)
	}

	b.translations[lang][key] = text
}
