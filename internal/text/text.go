// Package text provides utility functions for string manipulation.
package text

import "strings"

// TrimRightWhitespace removes trailing whitespace characters from a string.
func TrimRightWhitespace(s string) string {
	cutset := " \n\r\t"
	return strings.TrimRight(s, cutset)
}
