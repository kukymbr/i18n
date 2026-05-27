# i18n

[![License](https://img.shields.io/github/license/kukymbr/i18n.svg)](https://github.com/kukymbr/i18n/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/i18n.svg)](https://github.com/kukymbr/i18n/releases/latest)
[![GoDoc](https://godoc.org/github.com/kukymbr/i18n?status.svg)](https://godoc.org/github.com/kukymbr/i18n)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/i18n)](https://goreportcard.com/report/github.com/kukymbr/i18n)

Package to translate things.

This package helps you manage multi-language support in Go applications: loads translations from files and returns
localized strings with fallback logic.

## Installation

```bash
go get github.com/kukymbr/i18n
```

## Usage example

```go
package main

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
```

## Quick start

1. Create some YAML files with translations. 
   File could have any name, but you may like the next format: `<semantic_namespace>.<lang>.yaml`.
   Every file must contain the following structure:
   ```yaml
   # Code of the language presented in the file.
   language: en
   # Translations in key:value format.
   translations:
      # Keys could have any level of nesting.
      greeting: 
        hello: Hello!
   ```
   See the [testdata/example](testdata/example) for an example.
2. Create the bundle, using one of the `From*` functions (see [bundlesource.go](bundlesource.go)):
   ```go
   bundle, err := i18n.NewBundle(i18n.English, i18n.FromDirs(i18n.YAML, true, "testdata/example"))
   ```
3. Translate:
   ```go
   msg := bundle.Translate(i18n.English, "greeting.hello")
   ```

## Documentation

See the [Go reference](https://godoc.org/github.com/kukymbr/i18n).

## License

[MIT](LICENSE)
