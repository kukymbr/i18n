package i18n_test

import (
	"github.com/kukymbr/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDetectMissingTranslations(t *testing.T) {
	tests := []struct {
		Name   string
		Source i18n.BundleSource
		Assert func(t *testing.T, errs i18n.TranslationErrors)
	}{
		{
			Name:   "yaml bundle has errors",
			Source: i18n.FromEmbeddedFS(i18n.YAML, embeddedYAML, true, "testdata/yaml"),
			Assert: func(t *testing.T, errs i18n.TranslationErrors) {
				require.Len(t, errs, 6)
				assert.Equal(t, errs.Error(), `en: test_2: key is missing
es: errors.nested.test_6: key is missing
es: errors.nested.test_7: key is missing
es: errors.test_4: key is missing
es: errors.test_5: key is missing
es: test_4: value is empty`)
			},
		},
		{
			Name:   "json bundle has no errors",
			Source: i18n.FromEmbeddedFS(i18n.JSON, embeddedJSON, true, "testdata/json"),
			Assert: func(t *testing.T, errs i18n.TranslationErrors) {
				require.Empty(t, errs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			bundle, err := i18n.NewBundle(i18n.English, test.Source)
			require.NoError(t, err)
			require.NotNil(t, bundle)

			errs := i18n.DetectMissingTranslations(t.Context(), bundle)
			test.Assert(t, errs)
		})
	}
}
