package i18n_test

import (
	"testing"

	"github.com/kukymbr/i18n"
	"github.com/kukymbr/i18n/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				assert.Equal(t, "test", i18n.Translate(i18n.English, "test"))
			},
		},
		{
			Name: "when global bundle set",
			GetBundle: func(t *testing.T) *i18n.Bundle {
				b, err := i18n.NewBundle(
					i18n.English,
					i18n.FromFunc(func() (i18n.Tag, i18n.Translations, error) {
						return i18n.English, i18n.Translations{
							"test_1": "test text",
							"test_2": "test {{ .TestN }} text",
						}, nil
					}),
				)

				require.NoError(t, err)

				return b
			},
			Assert: func(t *testing.T) {
				tplData := tplData{TestN: 2}

				assert.Equal(t, "test text", i18n.Translate(i18n.English, "test_1"))
				assert.Equal(
					t,
					"test 2 text",
					i18n.T(i18n.English, "test_2", tplData),
				)

				ts := &testStruct{}
				err := i18n.TranslateStruct(i18n.English, ts, tplData)

				require.NoError(t, err)
				assert.Equal(t, testStruct{
					Test1: "test text",
					Test2: "test 2 text",
					Test3: "test_3",
				}, *ts)
			},
		},
		{
			Name: "with custom unmarshaler",
			GetBundle: func(t *testing.T) *i18n.Bundle {
				i18n.RegisterDataType("TEST", func(data []byte, v any) error {
					return json.Unmarshal([]byte(`{"language": "en", "translations": {"test": "test text"}}`), v)
				})

				b, err := i18n.NewBundle(i18n.English, i18n.FromString("TEST", "test"))

				require.NoError(t, err)

				return b
			},
			Assert: func(t *testing.T) {
				assert.Equal(t, "test text", i18n.Translate(i18n.English, "test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.GetBundle != nil {
				i18n.SetGlobalBundle(test.GetBundle(t))
			}

			require.NotPanics(t, func() {
				b := i18n.B()

				require.NotNil(t, b)
			})

			test.Assert(t)
		})
	}
}
