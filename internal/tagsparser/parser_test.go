package tagsparser_test

import (
	"strings"
	"testing"

	"github.com/kukymbr/i18n/internal/tagsparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name        string `i18n:"name.label"`
	Label       string
	Description string `i18n:"-"`
	Number      int
}

func TestParseTags(t *testing.T) {
	val := testStruct{
		Name:        "John Doe",
		Label:       "El es hombre",
		Description: "An example description",
		Number:      42,
	}

	require.NotPanics(t, func() {
		err := tagsparser.ParseTags(&val, func(s string) string {
			return strings.ReplaceAll(strings.ToLower(s), " ", "_")
		})

		require.NoError(t, err)
	})

	assert.Equal(t, "name.label", val.Name)
	assert.Equal(t, "el_es_hombre", val.Label)
	assert.Equal(t, "An example description", val.Description)
	assert.Equal(t, 42, val.Number)
}

func TestParseTags_NegativeCases(t *testing.T) {
	tests := []struct {
		Name         string
		GetInputFunc func() any
	}{
		{
			Name: "not a pointer",
			GetInputFunc: func() any {
				return testStruct{}
			},
		},
		{
			Name: "not a struct pointer",
			GetInputFunc: func() any {
				var v int

				return &v
			},
		},
		{
			Name: "not initialized struct",
			GetInputFunc: func() any {
				var v *testStruct

				return v
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			inp := test.GetInputFunc()

			require.NotPanics(t, func() {
				err := tagsparser.ParseTags(inp, func(s string) string {
					return s
				})

				require.Error(t, err)
			})
		})
	}
}

func TestParseTags_Panics(t *testing.T) {
	require.Panics(t, func() {
		_ = tagsparser.ParseTags(&testStruct{}, nil)
	})
}
