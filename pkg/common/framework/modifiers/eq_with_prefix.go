package modifiers

import "strings"

// EqualUpToPrefixFunc returns function that compares if two strings are equal stripping prefix
func EqualUpToPrefixFunc(prefix string) func(string, string) bool {
	return func(s1, s2 string) bool {
		return strings.TrimPrefix(s1, prefix) ==
			strings.TrimPrefix(s2, prefix)
	}
}
