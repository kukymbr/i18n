package i18n

import (
	"embed"
	"fmt"
	"io"
)

// BundleSource is a function adding Translations into the Bundle.
type BundleSource func(b *Bundle) error

// FromDirs reads Translations from the specified directories.
func FromDirs(dataType DataType, recursive bool, paths ...string) BundleSource {
	return func(b *Bundle) error {
		for _, path := range paths {
			if err := readFromDirectory(path, dataType, recursive, b.addTranslations); err != nil {
				return err
			}
		}

		return nil
	}
}

// FromFiles reads Translations from the specified files.
func FromFiles(dataType DataType, paths ...string) BundleSource {
	return func(b *Bundle) error {
		for _, path := range paths {
			lang, translations, err := readFromFile(path, dataType)
			if err != nil {
				return err
			}

			b.addTranslations(lang, translations)
		}

		return nil
	}
}

// FromEmbeddedFS reads Translations from the embed.FS.
func FromEmbeddedFS(dataType DataType, fs embed.FS, recursive bool, paths ...string) BundleSource {
	return func(b *Bundle) error {
		for _, path := range paths {
			err := readFromEmbeddedDirectory(fs, path, dataType, recursive, b.addTranslations)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// FromReader reads Translations from the specified io.Reader.
func FromReader(dataType DataType, r io.Reader) BundleSource {
	return func(b *Bundle) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to read from reader: %w", err)
		}

		return FromBytes(dataType, data)(b)
	}
}

// FromString parses Translations from the specified string.
func FromString(dataType DataType, inp string) BundleSource {
	return func(b *Bundle) error {
		return FromBytes(dataType, []byte(inp))(b)
	}
}

// FromBytes parses Translations from the specified bytes array.
func FromBytes(dataType DataType, inp []byte) BundleSource {
	return func(b *Bundle) error {
		lang, translations, err := readFromBytes(inp, dataType)
		if err != nil {
			return err
		}

		b.addTranslations(lang, translations)

		return nil
	}
}

// FromFunc adds Translations from the results of the callback fn.
func FromFunc(fn func() (Tag, Translations, error)) BundleSource {
	return func(b *Bundle) error {
		lang, translations, err := fn()
		if err != nil {
			return err
		}

		b.addTranslations(lang, translations)

		return nil
	}
}
