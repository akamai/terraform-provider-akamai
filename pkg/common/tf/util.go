package tf

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetToStringSlice converts schema.Set to a slice of strings
func SetToStringSlice(s *schema.Set) []string {
	list := make([]string, s.Len())
	for i, v := range s.List() {
		list[i] = v.(string)
	}
	return list
}

// ConvertListOfIntToInt64 casts slice of any type into slice of int64
func ConvertListOfIntToInt64(ints []interface{}) []int64 {
	var result []int64
	for _, v := range ints {
		result = append(result, int64(v.(int)))
	}
	return result
}

// MaxDuration returns the larger of x or y.
func MaxDuration(x, y time.Duration) time.Duration {
	if x < y {
		return y
	}
	return x
}
