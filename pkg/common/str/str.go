// Package str contains useful functions for string manipulation.
package str

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// From converts any type to string.
func From(data interface{}) string {
	var res string
	switch v := data.(type) {
	case float32, float64:
		res = fmt.Sprintf("%g", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		res = fmt.Sprintf("%d", v)
	case json.Number:
		res = v.String()
	case string:
		res = v
	case []byte:
		res = string(v)
	case bool:
		res = strconv.FormatBool(v)
	default:
		res = fmt.Sprintf("%v", data)
	}

	return res
}

// FirstNotEmpty returns first not empty string
func FirstNotEmpty(values ...string) string {
	for _, s := range values {
		if s != "" {
			return s
		}
	}
	return ""
}
