package i18n

import (
	"strings"

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
