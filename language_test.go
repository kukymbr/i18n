package i18n_test

import (
	"testing"

	"github.com/kukymbr/i18n"
	"github.com/stretchr/testify/assert"
)

func TestLang_Strings(t *testing.T) {
	tests := []struct {
		Input    string
		Fallback i18n.Tag
		Expected i18n.Tag
	}{
		{
			Input:    "en",
			Expected: i18n.English,
		},
		{
			Input:    "fr",
			Expected: i18n.French,
		},
		{
			Input:    "",
			Fallback: i18n.French,
			Expected: i18n.French,
		},
		{
			Input:    "unknown",
			Fallback: i18n.French,
			Expected: i18n.French,
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
		Input    i18n.Tag
		Expected i18n.Tag
	}{
		{
			Name:     "en",
			Input:    i18n.English,
			Expected: i18n.English,
		},
		{
			Name:     "fr",
			Input:    i18n.French,
			Expected: i18n.French,
		},
		{
			Name:     "und",
			Input:    i18n.Und,
			Expected: i18n.English,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			lang := i18n.Lang(test.Input)

			assert.Equal(t, test.Expected, lang)
		})
	}
}
