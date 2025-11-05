package i18n_test

import (
	"database/sql/driver"
	"testing"

	"github.com/kukymbr/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Input    string
		Expected i18n.Tag
	}{
		{Input: "en", Expected: i18n.English},
		{Input: "ru", Expected: i18n.Russian},
		{Input: "en-US", Expected: i18n.AmericanEnglish},
		{Input: "fr", Expected: i18n.French},
		{Input: "unknown", Expected: i18n.Und},
		{Input: "", Expected: i18n.Und},
	}

	for _, test := range tests {
		t.Run("Parse "+test.Input, func(t *testing.T) {
			tag, err := i18n.Parse(test.Input)

			if test.Expected != i18n.Und {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			assert.Equal(t, test.Expected.String(), tag.String())
		})

		t.Run("MustParse "+test.Input, func(t *testing.T) {
			var tag i18n.Tag

			fn := func() {
				tag = i18n.MustParse(test.Input)
			}

			if test.Expected != i18n.Und {
				require.NotPanics(t, fn)
			} else {
				require.Panics(t, fn)
			}

			assert.Equal(t, test.Expected.String(), tag.String())
		})
	}
}

func TestConvertJSON(t *testing.T) {
	t.Run("when valid", func(t *testing.T) {
		var tag i18n.Tag

		err := tag.UnmarshalJSON([]byte(`"en"`))
		require.NoError(t, err)

		assert.Equal(t, i18n.English, tag)

		b, err := tag.MarshalJSON()
		require.NoError(t, err)

		assert.Equal(t, []byte(`"en"`), b)
	})

	t.Run("when invalid", func(t *testing.T) {
		var tag i18n.Tag

		err := tag.UnmarshalJSON([]byte(`"invalid"`))
		assert.Equal(t, i18n.Und, tag)
		require.Error(t, err)
	})
}

func TestConvertText(t *testing.T) {
	t.Run("when valid", func(t *testing.T) {
		var tag i18n.Tag

		err := tag.UnmarshalText([]byte("en"))
		require.NoError(t, err)

		assert.Equal(t, i18n.English, tag)

		b, err := tag.MarshalText()
		require.NoError(t, err)

		assert.Equal(t, []byte("en"), b)
	})

	t.Run("when invalid", func(t *testing.T) {
		var tag i18n.Tag

		err := tag.UnmarshalText([]byte("invalid"))
		assert.Equal(t, i18n.Und, tag)
		require.Error(t, err)
	})
}

func TestScanValue(t *testing.T) {
	tests := []struct {
		Name          string
		Value         any
		ScanExpected  i18n.Tag
		ValueExpected driver.Value
	}{
		{
			Name:          "when valid",
			Value:         "en",
			ScanExpected:  i18n.English,
			ValueExpected: "en",
		},
		{
			Name:          "when invalid",
			Value:         "invalid",
			ScanExpected:  i18n.Und,
			ValueExpected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var tag i18n.Tag

			err := tag.Scan(test.Value)

			if test.ScanExpected != i18n.Und {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			val, err := tag.Value()
			require.NoError(t, err)

			assert.Equal(t, test.ScanExpected, tag)
			assert.Equal(t, test.ValueExpected, val)
		})
	}
}
