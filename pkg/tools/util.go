// Package tools is where some legacy provider functions were dropped
package tools

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

// MaxDuration returns the larger of x or y.
func MaxDuration(x, y time.Duration) time.Duration {
	if x < y {
		return y
	}
	return x
}