package i18n

import (
	"regexp"
	"sync"

	"github.com/kukymbr/i18n/json"
	"gopkg.in/yaml.v3"
)

// Input data types available by default.
const (
	YAML DataType = "YAML"
	JSON DataType = "JSON"
)

var dataTypeMu sync.RWMutex

var unmarshalers = map[DataType]UnmarshalerFunc{
	YAML: yaml.Unmarshal,
	JSON: json.Unmarshal,
}

var dataTypeFilters = map[DataType][]*regexp.Regexp{
	YAML: {regexp.MustCompile(`(?i)\.ya*ml$`)},
	JSON: {regexp.MustCompile(`(?i)\.json$`)},
}

// DataType is a bundle source data type.
type DataType string

// UnmarshalerFunc is a function to unmarshal data.
type UnmarshalerFunc func(data []byte, v any) error

// RegisterDataType registers or replaces the UnmarshalerFunc as an unmarshaler for a given DataType.
// If the fileNameFilters are given, them will be applied while filtering file names in directories.
// To remove an existing filters for a data type, use `nil` as third argument value:
// <code>
// i18n.RegisterDataType("YAML", yaml.Unmarshal, nil)
// </code>
func RegisterDataType(t DataType, fn UnmarshalerFunc, fileNameFilters ...*regexp.Regexp) {
	dataTypeMu.Lock()
	defer dataTypeMu.Unlock()

	unmarshalers[t] = fn

	if len(fileNameFilters) == 0 {
		return
	}

	if len(fileNameFilters) == 1 && fileNameFilters[0] == nil {
		delete(dataTypeFilters, t)

		return
	}

	dataTypeFilters[t] = fileNameFilters
}
