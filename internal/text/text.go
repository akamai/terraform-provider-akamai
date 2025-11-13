// Package text provides utility functions for string manipulation.
package text

import (
	"fmt"
	"slices"
	"strings"
)

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
	for _, part := range parts {
		res = append(res, strings.TrimSpace(part))
	}

	return res, nil
}
