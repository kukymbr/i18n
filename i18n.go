package i18n

import "golang.org/x/text/language"

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

// RegisterUnmarshaler sets the UnmarshalerFunc as an unmarshaler for a given data type.
func RegisterUnmarshaler(t DataType, fn UnmarshalerFunc) {
	unmarshalers[t] = fn
}

// Translate translates key using the global bundle.
func Translate[T Language](lang T, key string, tplData ...any) string {
	return GetGlobalBundle().Translate(Lang(lang), key, tplData...)
}

// T is a short alias of Translate function.
func T[T Language](lang T, key string, tplData ...any) string {
	return Translate[T](lang, key, tplData...)
}
