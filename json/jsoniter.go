//go:build jsoniter

package json

import jsoniter "github.com/json-iterator/go"

func Marshal(v any) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return jsoniter.Unmarshal(data, v)
}
