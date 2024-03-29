// Package date contains logic for handling operations on datetime values.
package date

import (
	"errors"
	"fmt"
	"time"
)

// ErrDateFormat is returned when there is an error parsing a string date.
var ErrDateFormat = errors.New("unable to parse date")

// DefaultFormat is the datetime format used across the provider.
const DefaultFormat = "2006-01-02T15:04:05Z"

// Parse parses the given string datetime using the default DefaultFormat format.
func Parse(value string) (time.Time, error) {
	return ParseFormat(DefaultFormat, value)
}

// ParseFormat parses the given string datetime using the provided format.
func ParseFormat(format, value string) (time.Time, error) {
	date, err := time.Parse(format, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrDateFormat, err.Error())
	}
	return date, nil
}
