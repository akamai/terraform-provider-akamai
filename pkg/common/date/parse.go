// Package date contains logic for handling operations on datetime values.
package date

import (
	"errors"
	"fmt"
	"time"
)

// ErrDateFormat is returned when there is an error parsing a string date.
var ErrDateFormat = errors.New("unable to parse date")

// ErrMarshal is returned when there is an error marshaling a time.Time date.
var ErrMarshal = errors.New("unable to marshal date")

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

// FormatRFC3339 returns a textual representation of time formatted according to the RFC3339 standard.
// RFC3339 is a subset of ISO 8601 producing the format "2006-01-02T15:04:05Z" (for a UTC time)
// which is commonly used in Akamai Open API.
func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatRFC3339Nano returns a textual representation of time formatted according to the RFC3339Nano standard.
// It returns an empty string if the provided time is equal to zero
func FormatRFC3339Nano(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339Nano)
}

// ToString returns given date in the string format
func ToString(value time.Time) (string, error) {
	bytes, err := value.MarshalText()
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrMarshal, err)
	}

	return string(bytes), nil
}

// CapDuration limits a given duration to a maximum value
func CapDuration(t, tMax time.Duration) time.Duration {
	if t > tMax {
		return tMax
	}
	return t
}
