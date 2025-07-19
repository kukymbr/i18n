package i18n

import (
	"bytes"
	"html/template"
	"strings"
	"sync"
)

var templateCache = struct {
	mu        sync.RWMutex
	templates map[string]*template.Template
}{
	templates: make(map[string]*template.Template),
}

func prepareText(key string, text string, tplData any) string {
	if strings.Contains(text, "{{") {
		tpl := getTemplate(key, text)
		if tpl == nil {
			return text
		}

		var b bytes.Buffer

		err := tpl.Execute(&b, tplData)
		if err != nil {
			return text
		}

		return b.String()
	}

	return text
}

func getTemplate(key string, text string) *template.Template {
	var err error

	templateCache.mu.Lock()
	defer templateCache.mu.Unlock()

	tpl, ok := templateCache.templates[key]
	if ok {
		return tpl
	}

	tpl = template.New(key)

	tpl, err = tpl.Parse(text)
	if err != nil {
		return nil
	}

	templateCache.templates[key] = tpl

	return tpl
}
