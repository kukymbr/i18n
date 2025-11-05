package i18n

import (
	"fmt"

	"golang.org/x/text/language"
)

type Language interface {
	~string | language.Tag | Tag
}

// Lang silently prepares the language.Tag for the input lang.
//
// If lang is a string, a language.Make function is called;
// if lang is a language.Tag, it is returned;
// if parsed or given tag is undefined, the fallback language will be returned.
func Lang[T Language](lang T, fallback ...Tag) Tag {
	fb := English
	if len(fallback) > 0 && fallback[0] != Und {
		fb = fallback[0]
	}

	switch val := any(lang).(type) {
	case Tag:
		if val == Und {
			return fb
		}

		return val
	case language.Tag:
		if val == language.Und {
			return fb
		}

		return Tag{val}
	default:
		s := fmt.Sprintf("%s", val)
		if s == "" {
			return fb
		}

		tag, _ := Parse(s)
		if tag != Und {
			return tag
		}
	}

	return fb
}
