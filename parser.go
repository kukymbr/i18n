package i18n

import (
	"fmt"
)

type unmarshalDTO struct {
	Language string `yaml:"language" json:"language" db:"language" bson:"language" xml:"language"`
	//nolint:lll
	Translations map[string]any `yaml:"translations" json:"translations" db:"translations" bson:"translations" xml:"translations"`
}

func unmarshal(dataType DataType, data []byte) (Tag, Translations, error) {
	dto := unmarshalDTO{}
	translations := Translations{}

	fn, err := getUnmarshaler(dataType)
	if err != nil {
		return Und, nil, err
	}

	if err := fn(data, &dto); err != nil {
		return Und, nil, fmt.Errorf("failed to unmarshal translations data: %w", err)
	}

	if err := parseTranslations("", dto.Translations, translations); err != nil {
		return Und, nil, fmt.Errorf("failed to parse translations: %w", err)
	}

	if dto.Language == "" {
		return Und, translations, nil
	}

	lang, err := Parse(dto.Language)
	if err != nil {
		return Und, nil, fmt.Errorf("failed to parse language '%s': %w", dto.Language, err)
	}

	return lang, translations, nil
}

func parseTranslations(parentKey string, inp map[string]any, target Translations) error {
	for k, v := range inp {
		key := k

		if parentKey != "" {
			key = parentKey + "." + k
		}

		if s, ok := v.(string); ok {
			target[key] = s

			continue
		}

		if m, ok := v.(map[string]any); ok {
			if err := parseTranslations(key, m, target); err != nil {
				return err
			}

			continue
		}

		return fmt.Errorf("key %s: expected string or map, got %T", key, v)
	}

	return nil
}

func getUnmarshaler(dataType DataType) (UnmarshalerFunc, error) {
	dataTypeMu.RLock()
	defer dataTypeMu.RUnlock()

	fn, ok := unmarshalers[dataType]
	if !ok {
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}

	return fn, nil
}
