package tagsparser

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	tagI18nKey = "i18n"
)

func ParseTags(inp any, updateValueFn func(string) string) error {
	if updateValueFn == nil {
		panic("updateValueFn must not be nil")
	}

	inpType := reflect.TypeOf(inp)
	if kind := inpType.Kind(); kind != reflect.Ptr {
		return fmt.Errorf("expected a pointer to a struct, got %s", kind)
	}

	elemType := inpType.Elem()
	if kind := elemType.Kind(); kind != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got pointer to a %s", kind)
	}

	inpValue := reflect.ValueOf(inp)
	if inpValue.IsNil() {
		return errors.New("got nil structure")
	}

	inpElem := inpValue.Elem()

	for i := range elemType.NumField() {
		field := elemType.Field(i)

		if field.Type.Kind() != reflect.String {
			continue
		}

		key := field.Tag.Get(tagI18nKey)
		if key == "-" {
			continue
		}

		valuePtr := inpElem.Field(i)
		value := valuePtr.String()

		if key == "" {
			key = value
		}

		valuePtr.SetString(updateValueFn(key))
	}

	return nil
}
