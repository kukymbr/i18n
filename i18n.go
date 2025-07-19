package i18n

import (
	"regexp"

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

// RegisterDataType registers or replaces the UnmarshalerFunc as an unmarshaler for a given DataType.
// If the fileNameFilters are given, them will be applied while filtering file names in directories.
// To remove an existing filters for a data type, use `nil` as third argument value:
// <code>
// i18n.RegisterDataType("YAML", yaml.Unmarshal, nil)
// </code>
func RegisterDataType(t DataType, fn UnmarshalerFunc, fileNameFilters ...*regexp.Regexp) {
	dataTypeMu.Lock()
	defer dataTypeMu.Unlock()

	unmarshalers[t] = fn

	if len(fileNameFilters) == 0 {
		return
	}

	if len(fileNameFilters) == 1 && fileNameFilters[0] == nil {
		delete(dataTypeFilters, t)

		return
	}

	dataTypeFilters[t] = fileNameFilters
}

// Translate translates key using the global bundle.
func Translate[T Language](lang T, key string, tplData ...any) string {
	return GetGlobalBundle().Translate(Lang(lang), key, tplData...)
}

// T is a short alias of Translate function.
func T[T Language](lang T, key string, tplData ...any) string {
	return Translate[T](lang, key, tplData...)
}
