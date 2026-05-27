package i18n_test

import (
	"embed"
	"fmt"
	"github.com/kukymbr/i18n"
)

//go:embed testdata/example/*
var translationsFS embed.FS

func TranslateThings() {
	// Load translations into a bundle.
	bundle, err := i18n.NewBundle(i18n.English, i18n.FromEmbeddedFS(i18n.YAML, translationsFS, true, "testdata/example"))
	if err != nil {
		panic(err)
	}

	// Retrieve translation by key and language.
	msg := bundle.Translate(i18n.English, "greeting.hello")
	fmt.Println(msg) // Hello!

	msg = bundle.Translate(i18n.Spanish, "greeting.hello_name", struct{ Name string }{"Mateo"})
	fmt.Println(msg) // ¡Hola, Mateo!

	// Fallback to the default language if missing.
	msg = bundle.Translate(i18n.French, "common.app_name")
	fmt.Println(msg) // i18n Example (from fallback en)
}
