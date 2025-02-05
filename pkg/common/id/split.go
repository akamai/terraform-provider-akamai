// Package id contains functions related to processing and formatting of id attributes for terraform resources and data sources
package id

import (
	"fmt"
	"strings"
)

// Split splits the provided id, separated by the ':' separator, into the expected number of parts.
func Split(id string, expectedNum int, example string) ([]string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != expectedNum {
		return nil, fmt.Errorf("id '%s' is incorrectly formatted: should be of form '%s'", id, example)
	}
	return parts, nil
}
