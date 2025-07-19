//go:build jsoniter

package json

import jsoniter "github.com/json-iterator/go"

func Unmarshal(data []byte, v any) error {
	return jsoniter.Unmarshal(data, v)
}
