package i18n_test

import (
	"testing"

	"github.com/kukymbr/i18n"
	"github.com/kukymbr/i18n/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestGlobalBundle(t *testing.T) {
	tests := []struct {
		Name      string
		GetBundle func(t *testing.T) *i18n.Bundle
		Assert    func(t *testing.T)
	}{
		{
			Name: "when no global bundle set",
			Assert: func(t *testing.T) {
				assert.Equal(t, "test", i18n.Translate(language.English, "test"))
			},
		},
		{
			Name: "when global bundle set",
			GetBundle: func(t *testing.T) *i18n.Bundle {
				b, err := i18n.NewBundle(
					language.English,
					i18n.FromFunc(func() (language.Tag, i18n.Translations, error) {
						return language.English, i18n.Translations{
							"test":  "test text",
							"test2": "test {{ .TestN }} text",
						}, nil
					}),
				)

				require.NoError(t, err)

				return b
			},
			Assert: func(t *testing.T) {
				assert.Equal(t, "test text", i18n.Translate(language.English, "test"))
				assert.Equal(
					t,
					"test 2 text",
					i18n.T(language.English, "test2", map[string]any{
						"TestN": "2",
					}),
				)
			},
		},
		{
			Name: "with custom unmarshaler",
			GetBundle: func(t *testing.T) *i18n.Bundle {
				i18n.RegisterUnmarshaler("TEST", func(data []byte, v any) error {
					return json.Unmarshal([]byte(`{"language": "en", "translations": {"test": "test text"}}`), v)
				})

				b, err := i18n.NewBundle(language.English, i18n.FromString("TEST", "test"))

				require.NoError(t, err)

				return b
			},
			Assert: func(t *testing.T) {
				assert.Equal(t, "test text", i18n.Translate(language.English, "test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.GetBundle != nil {
				i18n.SetGlobalBundle(test.GetBundle(t))
			}

			require.NotPanics(t, func() {
				b := i18n.GetGlobalBundle()

				require.NotNil(t, b)
			})

			test.Assert(t)
		})
	}
}
