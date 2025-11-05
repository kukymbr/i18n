package i18n

import (
	"golang.org/x/text/language"
)

var globalBundle *Bundle

// SetGlobalBundle sets a global Bundle instance.
func SetGlobalBundle(b *Bundle) {
	globalBundle = b
}

// GetGlobalBundle returns a global Bundle instance.
func GetGlobalBundle() *Bundle {
	if globalBundle == nil {
		var err error

		globalBundle, err = NewBundle(language.English)
		if err != nil {
			panic(err)
		}
	}

	return globalBundle
}

// B is a short alias for GetGlobalBundle function.
func B() *Bundle {
	return GetGlobalBundle()
}

// Translate translates key using the global bundle.
func Translate[T Language](lang T, key string, tplData ...any) string {
	return GetGlobalBundle().Translate(Lang(lang), key, tplData...)
}

// T is a short alias of Translate function.
func T[T Language](lang T, key string, tplData ...any) string {
	return Translate[T](lang, key, tplData...)
}

// TranslateStruct translates a structure using the global bundle.
// See Bundle.TranslateStruct for info.
func TranslateStruct(lang Tag, structure any, tplData ...any) error {
	return GetGlobalBundle().TranslateStruct(lang, structure, tplData...)
}
