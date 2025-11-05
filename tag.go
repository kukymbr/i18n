package i18n

import (
	"database/sql/driver"
	"fmt"

	"github.com/kukymbr/i18n/json"
	"golang.org/x/text/language"
)

var (
	Und = Tag{language.Und}

	Afrikaans            = Tag{language.Afrikaans}
	Amharic              = Tag{language.Amharic}
	Arabic               = Tag{language.Arabic}
	ModernStandardArabic = Tag{language.ModernStandardArabic}
	Azerbaijani          = Tag{language.Azerbaijani}
	Bulgarian            = Tag{language.Bulgarian}
	Bengali              = Tag{language.Bengali}
	Catalan              = Tag{language.Catalan}
	Czech                = Tag{language.Czech}
	Danish               = Tag{language.Danish}
	German               = Tag{language.German}
	Greek                = Tag{language.Greek}
	English              = Tag{language.English}
	AmericanEnglish      = Tag{language.AmericanEnglish}
	BritishEnglish       = Tag{language.BritishEnglish}
	Spanish              = Tag{language.Spanish}
	EuropeanSpanish      = Tag{language.EuropeanSpanish}
	LatinAmericanSpanish = Tag{language.LatinAmericanSpanish}
	Estonian             = Tag{language.Estonian}
	Persian              = Tag{language.Persian}
	Finnish              = Tag{language.Finnish}
	Filipino             = Tag{language.Filipino}
	French               = Tag{language.French}
	CanadianFrench       = Tag{language.CanadianFrench}
	Gujarati             = Tag{language.Gujarati}
	Hebrew               = Tag{language.Hebrew}
	Hindi                = Tag{language.Hindi}
	Croatian             = Tag{language.Croatian}
	Hungarian            = Tag{language.Hungarian}
	Armenian             = Tag{language.Armenian}
	Indonesian           = Tag{language.Indonesian}
	Icelandic            = Tag{language.Icelandic}
	Italian              = Tag{language.Italian}
	Japanese             = Tag{language.Japanese}
	Georgian             = Tag{language.Georgian}
	Kazakh               = Tag{language.Kazakh}
	Khmer                = Tag{language.Khmer}
	Kannada              = Tag{language.Kannada}
	Korean               = Tag{language.Korean}
	Kirghiz              = Tag{language.Kirghiz}
	Lao                  = Tag{language.Lao}
	Lithuanian           = Tag{language.Lithuanian}
	Latvian              = Tag{language.Latvian}
	Macedonian           = Tag{language.Macedonian}
	Malayalam            = Tag{language.Malayalam}
	Mongolian            = Tag{language.Mongolian}
	Marathi              = Tag{language.Marathi}
	Malay                = Tag{language.Malay}
	Burmese              = Tag{language.Burmese}
	Nepali               = Tag{language.Nepali}
	Dutch                = Tag{language.Dutch}
	Norwegian            = Tag{language.Norwegian}
	Punjabi              = Tag{language.Punjabi}
	Polish               = Tag{language.Polish}
	Portuguese           = Tag{language.Portuguese}
	BrazilianPortuguese  = Tag{language.BrazilianPortuguese}
	EuropeanPortuguese   = Tag{language.EuropeanPortuguese}
	Romanian             = Tag{language.Romanian}
	Russian              = Tag{language.Russian}
	Sinhala              = Tag{language.Sinhala}
	Slovak               = Tag{language.Slovak}
	Slovenian            = Tag{language.Slovenian}
	Albanian             = Tag{language.Albanian}
	Serbian              = Tag{language.Serbian}
	SerbianLatin         = Tag{language.SerbianLatin}
	Swedish              = Tag{language.Swedish}
	Swahili              = Tag{language.Swahili}
	Tamil                = Tag{language.Tamil}
	Telugu               = Tag{language.Telugu}
	Thai                 = Tag{language.Thai}
	Turkish              = Tag{language.Turkish}
	Ukrainian            = Tag{language.Ukrainian}
	Urdu                 = Tag{language.Urdu}
	Uzbek                = Tag{language.Uzbek}
	Vietnamese           = Tag{language.Vietnamese}
	Chinese              = Tag{language.Chinese}
	SimplifiedChinese    = Tag{language.SimplifiedChinese}
	TraditionalChinese   = Tag{language.TraditionalChinese}
	Zulu                 = Tag{language.Zulu}
)

type Tag struct {
	language.Tag
}

func Parse[T ~string](s T) (Tag, error) {
	tag, err := language.Parse(string(s))
	if err != nil {
		return Und, err
	}

	return Tag{tag}, nil
}

func MustParse[T ~string](s T) Tag {
	tag, err := Parse(s)
	if err != nil {
		panic(err)
	}

	return tag
}

func (t Tag) String() string {
	if t == Und {
		return ""
	}

	return t.Tag.String()
}

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Tag) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("unmarshal JSON into language tag string: %w", err)
	}

	tag, err := Parse(s)
	if err != nil {
		return fmt.Errorf("parse language tag string from JSON: %w", err)
	}

	*t = tag

	return nil
}

func (t Tag) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Tag) UnmarshalText(data []byte) error {
	tag, err := Parse(string(data))
	if err != nil {
		return fmt.Errorf("parse language tag string from text: %w", err)
	}

	*t = tag

	return nil
}

func (t *Tag) Scan(value any) error {
	if value == nil {
		*t = Und

		return nil
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("scan i18n tag value: expected string, got %T", value)
	}

	tag, err := Parse(s)
	if err != nil {
		return err
	}

	*t = tag

	return nil
}

func (t *Tag) Value() (driver.Value, error) {
	return t.String(), nil
}
