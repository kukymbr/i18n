package i18n

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func readFromDirectory(path string, dataType DataType, recursive bool, each func(Tag, Translations)) error {
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

		if !acceptFile(dataType, entry.Name()) {
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
	each func(Tag, Translations),
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

			err := readFromEmbeddedDirectory(fs, entryPath, dataType, recursive, each)
			if err != nil {
				return fmt.Errorf("%s: %w", entryPath, err)
			}

			continue
		}

		if !acceptFile(dataType, entry.Name()) {
			return nil
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
	each func(Tag, Translations),
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

func readFromFile(path string, dataType DataType) (Tag, Translations, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Und, nil, fmt.Errorf("failed to read i18n file '%s': %w", path, err)
	}

	return readFromBytes(data, dataType)
}

func readFromBytes(data []byte, dataType DataType) (Tag, Translations, error) {
	return unmarshal(dataType, data)
}

func acceptFile(dataType DataType, name string) bool {
	dataTypeMu.RLock()
	defer dataTypeMu.RUnlock()

	rxs, ok := dataTypeFilters[dataType]
	if !ok {
		return true
	}

	for _, rx := range rxs {
		if rx.MatchString(name) {
			return true
		}
	}

	return false
}
