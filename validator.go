package i18n

import (
	"context"
	"errors"
	"slices"
	"strings"
)

var (
	ErrTranslationKeyMissing = errors.New("key is missing")
	ErrTranslationKeyEmpty   = errors.New("value is empty")
)

// DetectMissingTranslations checks the Bundle and reports all the missing keys.
func DetectMissingTranslations(ctx context.Context, bundle *Bundle) TranslationErrors {
	allKeys := make(map[string]bool)
	index := make(map[Tag]map[string]bool, len(bundle.translations))

	for tag, translations := range bundle.translations {
		if ctx.Err() != nil {
			//nolint:nilerr // FP.
			return nil
		}

		if index[tag] == nil {
			index[tag] = make(map[string]bool, len(translations))
		}

		for key, val := range translations {
			hasValue := strings.TrimSpace(val) != ""
			flagVal := hasValue

			if alreadyHasValue, ok := allKeys[key]; ok && alreadyHasValue {
				flagVal = alreadyHasValue
			}

			allKeys[key] = flagVal
			index[tag][key] = hasValue
		}
	}

	var errs TranslationErrors

	for key, requiredValue := range allKeys {
		if ctx.Err() != nil {
			//nolint:nilerr // FP.
			return nil
		}

		for tag, keys := range index {
			hasValue, ok := keys[key]
			if !ok {
				errs = append(errs, NewTranslationError(tag, key, ErrTranslationKeyMissing))

				continue
			}

			if !hasValue && requiredValue {
				errs = append(errs, NewTranslationError(tag, key, ErrTranslationKeyEmpty))
			}
		}
	}

	slices.SortFunc(errs, func(a TranslationError, b TranslationError) int {
		if v := strings.Compare(a.Tag.String(), b.Tag.String()); v != 0 {
			return v
		}

		return strings.Compare(a.Key, b.Key)
	})

	return errs
}

// TranslationErrors is a list of TranslationError.
type TranslationErrors []TranslationError

func (errs TranslationErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}

	msg := strings.Builder{}

	for i, e := range errs {
		msg.WriteString(e.Error())

		if i < len(errs)-1 {
			msg.WriteString("\n")
		}
	}

	return msg.String()
}

// TranslationError is an error pointing to the missing key or empty value in the Tag's translations in the Bundle.
type TranslationError struct {
	Tag Tag
	Key string
	Err error
}

func NewTranslationError(tag Tag, key string, err error) TranslationError {
	return TranslationError{tag, key, err}
}

func (e TranslationError) Error() string {
	return e.Tag.String() + ": " + e.Key + ": " + e.Err.Error()
}

func (e TranslationError) Unwrap() error {
	return e.Err
}
