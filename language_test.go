package i18n_test

import (
	"testing"

	"github.com/kukymbr/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLang_Strings(t *testing.T) {
	tests := []struct {
		Input    string
		Fallback language.Tag
		Expected language.Tag
	}{
		{
			Input:    "en",
			Expected: language.English,
		},
		{
			Input:    "fr",
			Expected: language.French,
		},
		{
			Input:    "",
			Fallback: language.French,
			Expected: language.French,
		},
		{
			Input:    "unknown",
			Fallback: language.French,
			Expected: language.French,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			lang := i18n.Lang(test.Input, test.Fallback)

			assert.Equal(t, test.Expected, lang)
		})
	}
}

func TestLang_Tags(t *testing.T) {
	tests := []struct {
		Name     string
		Input    language.Tag
		Expected language.Tag
	}{
		{
			Name:     "en",
			Input:    language.English,
			Expected: language.English,
		},
		{
			Name:     "fr",
			Input:    language.French,
			Expected: language.French,
		},
		{
			Name:     "und",
			Input:    language.Und,
			Expected: language.English,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			lang := i18n.Lang(test.Input)

			assert.Equal(t, test.Expected, lang)
		})
	}
}
