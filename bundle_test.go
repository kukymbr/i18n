package i18n_test

import (
	"embed"
	"errors"
	"fmt"
	"strings"
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

type structTranslationAssertion struct {
	Lang     language.Tag
	TplData  *tplData
	Expected testStruct
}

func TestBundle(t *testing.T) {
	tests := []struct {
		Name     string
		Fallback language.Tag
		Sources  []i18n.BundleSource

		AssertNew       func(t *testing.T, b *i18n.Bundle, err error)
		AssertTranslate []translationAssertion

		InputTranslateStruct  *testStruct
		AssertTranslateStruct []structTranslationAssertion
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
			InputTranslateStruct: &testStruct{
				TestSkip:       "test_1",
				TestWithoutTag: "test_1",
			},
			AssertTranslateStruct: []structTranslationAssertion{
				{
					Lang:    language.English,
					TplData: &tplData{TestN: 3},
					Expected: testStruct{
						Test1:          "Test 1 in JSON",
						Test2:          "Test 2 in JSON",
						Test3:          "Test 3 in JSON",
						TestSkip:       "test_1",
						TestWithoutTag: "Test 1 in JSON",
					},
				},
			},
		},
		{
			Name:     "From YAML dir (nested keys)",
			Fallback: language.English,
			Sources: []i18n.BundleSource{
				i18n.FromDirs(i18n.YAML, false, "testdata/yaml"),
			},
			AssertNew: func(t *testing.T, _ *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Key: "test_1", Expected: "Test 1 in YAML"},
				{Key: "errors.test_4", Expected: "Error 1"},
				{Key: "errors.test_5", Expected: "Error 2"},
				{Key: "errors.nested.test_6", Expected: "Error 3"},
				{Key: "errors.nested.test_7", Expected: "Error 4"},
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
				i18n.FromReader(i18n.JSON, strings.NewReader(`{"language": "it", "translations": {"test_1": "Test 1 del lettore"}}`)),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.NoError(t, err)
			},
			AssertTranslate: []translationAssertion{
				{Lang: language.English, Key: "test_1", Expected: "Test 1 in YAML"},
				{Lang: language.Spanish, Key: "test_1", Expected: "Prueba 1 en JSON"},
				{Lang: language.Russian, Key: "test_1", Expected: "Тест 1 из строки"},
				{Lang: language.Hebrew, Key: "test_1", Expected: "בדיקה 1 מבתים"},
				{Lang: language.Italian, Key: "test_1", Expected: "Test 1 del lettore"},
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
		{
			Name: "When data type is unknown",
			Sources: []i18n.BundleSource{
				i18n.FromString("UNKNOWN", "some input"),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When invalid translation type",
			Sources: []i18n.BundleSource{
				i18n.FromString(i18n.JSON, `{"translations": {"test": 0}}`),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When nested translation is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromString(i18n.JSON, `{"translations": {"nested": {"test": 0}}}`),
			},
			AssertNew: func(t *testing.T, b *i18n.Bundle, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "When language is invalid",
			Sources: []i18n.BundleSource{
				i18n.FromString(i18n.JSON, `{"language": "invalid"}`),
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

					text := bundle.T(tt.Lang, tt.Key, tplData...)

					assert.Equal(t, tt.Expected, text)
				})
			}

			if test.InputTranslateStruct == nil {
				return
			}

			for i, tt := range test.AssertTranslateStruct {
				t.Run(test.Name+":struct:"+fmt.Sprintf("%d", i), func(t *testing.T) {
					var (
						inp     = &testStruct{}
						tplData []any
					)

					if tt.TplData != nil {
						tplData = append(tplData, tt.TplData)
					}

					*inp = *test.InputTranslateStruct

					err := bundle.TranslateStruct(tt.Lang, inp, tplData...)
					require.NoError(t, err)

					assert.Equal(t, tt.Expected, *inp)
				})
			}
		})
	}
}

func TestEmptyBundle(t *testing.T) {
	bundle := i18n.NewEmptyBundle()

	assert.NotPanics(t, func() {
		assert.Equal(t, "Test 1", bundle.T(language.Japanese, "Test 1"))
		assert.Equal(t, "Test 2", bundle.Translate(language.Hungarian, "Test 2"))

		test3 := testStruct{
			Test1:          "Test1",
			Test2:          "Test2",
			Test3:          "Test3",
			TestSkip:       "TestSkip",
			TestWithoutTag: "TestWithoutTag",
		}

		require.NoError(t, bundle.TranslateStruct(language.Amharic, &test3))

		assert.Equal(t, "test_1", test3.Test1)
		assert.Equal(t, "test_2", test3.Test2)
		assert.Equal(t, "test_3", test3.Test3)
		assert.Equal(t, "TestSkip", test3.TestSkip)
		assert.Equal(t, "TestWithoutTag", test3.TestWithoutTag)
	})
}
