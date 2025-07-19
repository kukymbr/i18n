package i18n_test

import (
	"embed"
	"errors"
	"testing"

	"github.com/kukymbr/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

//go:embed all:testdata/json
var embeddedJSON embed.FS

//go:embed testdata/json_invalid
var embeddedInvalidJSON embed.FS

type translationAssertion struct {
	Lang     language.Tag
	Key      string
	TplData  *tplData
	Expected string
}

type tplData struct {
	TestN int
}

func TestBundle(t *testing.T) {
	tests := []struct {
		Name            string
		Fallback        language.Tag
		Sources         []i18n.BundleSource
		AssertNew       func(t *testing.T, b *i18n.Bundle, err error)
		AssertTranslate []translationAssertion
	}{
		// region Positive cases
		{
			Name:     "From JSON recursive dirs",
			Fallback: language.English,
			Sources: []i18n.BundleSource{
				i18n.FromDirs(i18n.JSON, true, "testdata/json"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Lang: language.English, Key: "test_1", Expected: "Test 1 in JSON"},
				{Lang: language.English, TplData: &tplData{TestN: 5}, Key: "test_3", Expected: "Test 5 in JSON"},
				{Lang: language.English, Key: "not from translations", Expected: "not from translations"},
				{Lang: language.Spanish, Key: "test_1", Expected: "Prueba 1 en JSON"},
				{Key: "test_1", Expected: "Test 1 in JSON"},
			},
		},
		{
			Name:     "From JSON non-recursive dirs",
			Fallback: language.English,
			Sources: []i18n.BundleSource{
				i18n.FromDirs(i18n.JSON, false, "testdata/json"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Lang: language.English, Key: "test_1", Expected: "Test 1 in JSON"},
				{Lang: language.Spanish, Key: "test_1", Expected: "Test 1 in JSON"},
			},
		},
		{
			Name:     "From multiple sources",
			Fallback: language.Russian,
			Sources: []i18n.BundleSource{
				i18n.FromFiles(i18n.YAML, "testdata/yaml/en.yml"),
				i18n.FromFiles(i18n.JSON, "testdata/json/es/es.json"),
				i18n.FromString(i18n.JSON, `{"translations": {"test_1": "Тест 1 из строки"}}`),
				i18n.FromBytes(i18n.JSON, []byte(`{"language": "he", "translations": {"test_1": "בדיקה 1 מבתים"}}`)),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Lang: language.English, Key: "test_1", Expected: "Test 1 in YAML"},
				{Lang: language.Spanish, Key: "test_1", Expected: "Prueba 1 en JSON"},
				{Lang: language.Russian, Key: "test_1", Expected: "Тест 1 из строки"},
				{Lang: language.Hebrew, Key: "test_1", Expected: "בדיקה 1 מבתים"},
			},
		},
		{
			Name: "From func",
			Sources: []i18n.BundleSource{
				i18n.FromFunc(func() (language.Tag, i18n.Translations, error) {
					return language.Italian, i18n.Translations{"test_1": "Test 1 da callback"}, nil
				}),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Lang: language.English, Key: "test_1", Expected: "test_1"},
				{Lang: language.Italian, Key: "test_1", Expected: "Test 1 da callback"},
			},
		},
		{
			Name: "From embedded JSON",
			Sources: []i18n.BundleSource{
				i18n.FromEmbeddedFS(i18n.JSON, embeddedJSON, true, "testdata/json"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Lang: language.English, Key: "test_1", Expected: "Test 1 in JSON"},
				{Lang: language.Spanish, Key: "test_1", Expected: "Prueba 1 en JSON"},
			},
		},
		// endregion Positive cases

		// region Negative cases
		{
			Name: "When dir is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromDirs(i18n.JSON, false, "testdata/unknown"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When content in file is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromFiles(i18n.JSON, "testdata/json_invalid/invalid.json"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When content in dir is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromDirs(i18n.JSON, true, "testdata/json_invalid"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When embedded JSON is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromEmbeddedFS(i18n.JSON, embeddedInvalidJSON, true, "testdata/json"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When callback returns error",
			Sources: []i18n.BundleSource{
				i18n.FromFunc(func() (language.Tag, i18n.Translations, error) {
					return language.Und, nil, errors.New("test error")
				}),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When string is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromString(i18n.JSON, "{{ some broken JSON }}"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		// endregion Negative cases
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			bundle, err := i18n.NewBundle(test.Fallback, test.Sources...)

			test.AssertNew(t, bundle, err)

			for _, tt := range test.AssertTranslate {
				t.Run(test.Name+":"+tt.Key, func(t *testing.T) {
					var tplData []any
					if tt.TplData != nil {
						tplData = append(tplData, tt.TplData)
					}

					text := bundle.Translate(tt.Lang, tt.Key, tplData...)

					assert.Equal(t, tt.Expected, text)
				})
			}
		})
	}
}
