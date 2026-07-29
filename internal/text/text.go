// Package text provides utility functions for string manipulation.
package text

import (
	"fmt"
	"slices"
	"strings"
)

// JoinStringBased is a generic version of strings.Join that works with any type based on string.
// It converts the elements to string before joining them with the given separator.
func JoinStringBased[T ~string](elems []T, sep string) string {
	ss := make([]string, len(elems))
	for i, v := range elems {
		ss[i] = string(v)
	}
	return strings.Join(ss, sep)
}

// ToStrings converts a slice of elements of type T (where T is a type whose underlying type is string)
// to a slice of strings. Each element is converted to a string and placed in the resulting slice.
func ToStrings[T ~string](elems []T) []string {
	ss := make([]string, len(elems))
	for i, v := range elems {
		ss[i] = string(v)
	}
	return ss
}

// TrimRightWhitespace removes trailing whitespace characters from a string.
func TrimRightWhitespace(s string) string {
	cutset := " \n\r\t"
	return strings.TrimRight(s, cutset)
}

// IDSplitter helps to split and validate import IDs.
type IDSplitter struct {
	FormatHint      string
	AcceptedLengths []int
}

// ImportIDSplitter creates a new IDSplitter with a format hint.
func ImportIDSplitter(formatHint string) IDSplitter {
	return IDSplitter{
		FormatHint: formatHint,
	}
}

// AcceptLen adds an accepted length for the split ID parts.
func (p IDSplitter) AcceptLen(length int) IDSplitter {
	p.AcceptedLengths = append(p.AcceptedLengths, length)
	return p
}

// Split splits the given ID by commas and validates the number of parts.
func (p IDSplitter) Split(id string) ([]string, error) {
	if p.FormatHint == "" {
		return nil, fmt.Errorf("no format hint defined for importID; you need to provide a format hint using ImportIDSplitter method")
	}

	if len(p.AcceptedLengths) == 0 {
		return nil, fmt.Errorf("no accepted lengths defined for importID; you need to provide at least one accepted length using AcceptLen method")
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("importID cannot be empty; you need to provide an importID in the format '%s'",
			p.FormatHint)
	}

	parts := strings.Split(id, ",")

	if !slices.Contains(p.AcceptedLengths, len(parts)) {
		return nil, fmt.Errorf("invalid number of importID parts: %d; you need to provide an importID in the format '%s'",
			len(parts), p.FormatHint)
	}

	var res []string
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("importID part %d cannot be empty; you need to provide an importID in the format '%s'",
				i+1, p.FormatHint)
		}
		res = append(res, part)
	}

	return res, nil
}
