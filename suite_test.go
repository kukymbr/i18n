package i18n_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

type tplData struct {
	TestN int
}

type testStruct struct {
	Test1          string `i18n:"test_1"`
	Test2          string `i18n:"test_2"`
	Test3          string `i18n:"test_3"`
	TestSkip       string `i18n:"-"`
	TestWithoutTag string
}

func TestReadme(t *testing.T) {
	originStdout := os.Stdout
	t.Cleanup(func() {
		os.Stdout = originStdout
	})

	r, w, _ := os.Pipe()
	os.Stdout = w

	require.NotPanics(t, TranslateThings)

	_ = w.Close()
	os.Stdout = originStdout

	recorded, err := io.ReadAll(r)
	require.NoError(t, err)

	assert.Equal(
		t,
		"Hello!\n"+
			"¡Hola, Mateo!\n"+
			"i18n Example\n",
		string(recorded),
	)
}
