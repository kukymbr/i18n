package i18n

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/kukymbr/i18n/json"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// Input data types available by default.
const (
	YAML DataType = "YAML"
	JSON DataType = "JSON"
)

var unmarshalers = map[DataType]UnmarshalerFunc{
	YAML: yaml.Unmarshal,
	JSON: json.Unmarshal,
}

// DataType is a bundle source data type.
type DataType string

// UnmarshalerFunc is a function to unmarshal data.
type UnmarshalerFunc func(data []byte, v any) error

type bundleDTO struct {
	Language string `yaml:"language" json:"language" db:"language" bson:"language" xml:"language"`
	//nolint:lll
	Translations Translations `yaml:"translations" json:"translations" db:"translations" bson:"translations" xml:"translations"`
}

func readFromDirectory(
	path string,
	dataType DataType,
	recursive bool,
	each func(language.Tag, Translations),
) error {
	err := filepath.WalkDir(path, func(entryPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == entryPath {
			return nil
		}

		if strings.HasPrefix(entry.Name(), ".") {
			if entry.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if entry.IsDir() {
			if !recursive {
				return filepath.SkipDir
			}

			return nil
		}

		lang, translations, err := readFromFile(entryPath, dataType)
		if err != nil {
			return err
		}

		each(lang, translations)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func readFromEmbeddedDirectory(
	fs embed.FS,
	path string,
	dataType DataType,
	recursive bool,
	each func(language.Tag, Translations),
) error {
	entries, err := fs.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read embedded directory: %w", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			if !recursive {
				continue
			}

			err := readFromEmbeddedDirectory(fs, filepath.Join(entryPath), dataType, recursive, each)
			if err != nil {
				return fmt.Errorf("%s: %w", entryPath, err)
			}

			continue
		}

		if err := readFromEmbeddedFile(fs, entryPath, dataType, each); err != nil {
			return err
		}
	}

	return nil
}

func readFromEmbeddedFile(
	fs embed.FS,
	path string,
	dataType DataType,
	each func(language.Tag, Translations),
) error {
	data, err := fs.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %w", path, err)
	}

	lang, translations, err := readFromBytes(data, dataType)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	each(lang, translations)

	return nil
}

func readFromFile(path string, dataType DataType) (language.Tag, Translations, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return language.Tag{}, nil, fmt.Errorf("failed to read i18n file '%s': %w", path, err)
	}

	return readFromBytes(data, dataType)
}

func readFromBytes(data []byte, dataType DataType) (language.Tag, Translations, error) {
	fn, ok := unmarshalers[dataType]
	if !ok {
		return language.Tag{}, nil, fmt.Errorf("unsupported data type: %s", dataType)
	}

	dto := &bundleDTO{}

	if err := fn(data, &dto); err != nil {
		return language.Tag{}, nil, fmt.Errorf("failed to unmarshal translations data: %w", err)
	}

	if dto.Language == "" {
		return language.Und, dto.Translations, nil
	}

	lang, err := language.Parse(dto.Language)
	if err != nil {
		return language.Tag{}, nil, fmt.Errorf("failed to parse language '%s': %w", dto.Language, err)
	}

	return lang, dto.Translations, nil
}
