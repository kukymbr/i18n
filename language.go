package i18n

import (
	"fmt"

	"golang.org/x/text/language"
)

type Language interface {
	~string | language.Tag
}

// Lang silently prepares the language.Tag for the input lang.
//
// If lang is a string, a language.Make function is called;
// if lang is a language.Tag, it is returned;
// if parsed or given tag is undefined, the fallback language will be returned.
func Lang[T Language](lang T, fallback ...language.Tag) language.Tag {
	fb := language.English

	if len(fallback) > 0 && fallback[0] != language.Und {
		fb = fallback[0]
	}

	switch val := any(lang).(type) {
	case language.Tag:
		if val == language.Und {
			return fb
		}

		return val
	default:
		s := fmt.Sprintf("%s", val)
		if s == "" {
			return fb
		}

		tag := language.Make(s)
		if tag != language.Und {
			return tag
		}
	}

	return fb
}
