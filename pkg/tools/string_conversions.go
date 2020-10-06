package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ConvertToString(data interface{}) (res string) {
	switch v := data.(type) {
	case float64:
		res = fmt.Sprint(data.(float64))
	case float32:
		res = fmt.Sprint(data.(float32))
	case int, int64:
		res = strconv.Itoa(data.(int))
	case json.Number:
		res = data.(json.Number).String()
	case string:
		res = data.(string)
	case []byte:
		res = string(v)
	default:
		res = ""
	}
	return
}
