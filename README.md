# i18n

[![License](https://img.shields.io/github/license/kukymbr/i18n.svg)](https://github.com/kukymbr/i18n/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/i18n.svg)](https://github.com/kukymbr/i18n/releases/latest)
[![GoDoc](https://godoc.org/github.com/kukymbr/i18n?status.svg)](https://godoc.org/github.com/kukymbr/i18n)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/i18n)](https://goreportcard.com/report/github.com/kukymbr/i18n)

Package to translate things.

This package helps you manage multi-language support in Go applications — from loading translation files to retrieving
localized strings with fallback logic.

## Installation

```bash
go get github.com/kukymbr/i18n
```

## Quick Start

```go
package main

import (
	"embed"
	"fmt"

	"github.com/kukymbr/i18n"
)

//go:embed translations/*
var i18nFS embed.FS

func main() {
	// Load translations into a bundle
	bundle, err := i18n.NewBundle("en", i18n.FromEmbeddedFS(i18n.YAML, i18nFS, true))
	if err != nil {
	    panic(err)
	}

	// Retrieve translation by key and language
	msg := bundle.T(i18n.English, "greeting.hello")
	fmt.Println(msg) // Hello

	msg = bundle.T(i18n.Spanish, "greeting.hello")
	fmt.Println(msg) // Hola

	// Fallback to the default language if missing
	msg = bundle.T(i18n.French, "greeting.hello")
	fmt.Println(msg) // Hello (from default en)
}
```

## License

[MIT](LICENSE)
