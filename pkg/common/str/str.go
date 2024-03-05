// Package str contains useful functions for string manipulation.
package str

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// From converts any type to string.
func From(data interface{}) string {
	switch v := data.(type) {
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case json.Number:
		return v.String()
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", data)
	}
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
