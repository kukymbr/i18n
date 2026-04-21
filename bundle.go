package i18n

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/kukymbr/i18n/internal/tagsparser"
)

// Translations is a map of translations in a key:text format
type Translations map[string]string

// Bundle is an i18n translations bundle.
type Bundle struct {
	fallbackLanguage Tag
	translations     map[Tag]Translations
}

// NewBundle creates a new Bundle instance.
func NewBundle(fallbackLanguage Tag, sources ...BundleSource) (*Bundle, error) {
	b := &Bundle{fallbackLanguage: fallbackLanguage}

	b.translations = make(map[Tag]Translations)

	for _, source := range sources {
		if err := source(b); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// NewEmptyBundle returns a new Bundle without any translations.
func NewEmptyBundle() *Bundle {
	b, _ := NewBundle(English)

	return b
}

// Translate finds a translation for a key.
func (b *Bundle) Translate(lang Tag, key string, tplData ...any) string {
	var data any
	if len(tplData) > 0 {
		data = tplData[0]
	}

	return b.translate(lang, key, data)
}

// T is a short alias for a Translate.
func (b *Bundle) T(lang Tag, key string, tplData ...any) string {
	return b.Translate(lang, key, tplData...)
}

// TranslateStruct updated fields of the given structure with a translated representation.
// The structure argument must be a pointer to a non-nil structure variable.
//
// The `i18n:"field.key"` tag format is expected to get a field's translation key;
// if no `i18n ` is found, the field's value is used as a key.
// Add `i18n:"-"` tag to skip field's translation.
// Only string values are affected.
func (b *Bundle) TranslateStruct(lang Tag, structure any, tplData ...any) error {
	err := tagsparser.ParseTags(structure, func(s string) string {
		return b.Translate(lang, s, tplData...)
	})
	if err != nil {
		return fmt.Errorf("translate structure: %w", err)
	}

	return nil
}

// GetFallbackLanguage returns the fallback language.
func (b *Bundle) GetFallbackLanguage() Tag {
	return b.fallbackLanguage
}

// CalcHash calculates hash of the whole bundle.
func (b *Bundle) CalcHash() (string, error) {
	// Prepare sorted list of tags to avoid random hash changes because of unstable map keys order.
	tags := getSortedKeys(b.translations, func(a Tag, b Tag) int {
		return strings.Compare(a.String(), b.String())
	})

	hasher := sha256.New()

	hasher.Write([]byte("_fallback:" + b.fallbackLanguage.String() + ";"))

	for _, tag := range tags {
		keys := getSortedKeys(b.translations[tag], strings.Compare)

		for _, key := range keys {
			hasher.Write([]byte(key + ":" + b.translations[tag][key] + ";"))
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Translate finds a translation for a key.
func (b *Bundle) translate(lang Tag, key string, tplData any) string {
	if lang == Und {
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

func (b *Bundle) getTranslation(lang Tag, key string) (string, bool) {
	text, ok := b.translations[lang][key]

	return text, ok
}

func (b *Bundle) addTranslations(lang Tag, translations Translations) {
	for key, text := range translations {
		b.addTranslation(lang, key, text)
	}
}

func (b *Bundle) addTranslation(lang Tag, key string, text string) {
	if lang == Und {
		lang = b.fallbackLanguage
	}

	if _, ok := b.translations[lang]; !ok {
		b.translations[lang] = make(Translations)
	}

	b.translations[lang][key] = text
}

func getSortedKeys[Tk comparable, Tv any](s map[Tk]Tv, compareFunc func(a Tk, b Tk) int) []Tk {
	keys := make([]Tk, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}

	slices.SortStableFunc(keys, compareFunc)

	return keys
}
