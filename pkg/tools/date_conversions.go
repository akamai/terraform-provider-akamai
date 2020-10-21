package tools

import (
	"errors"
	"fmt"
	"time"
)

// ErrDateFormat is returned when there is an error parsing a string date
var ErrDateFormat = errors.New("unable to parse date")

const DateTimeFormat = "2006-01-02T15:04:05Z"

func ParseDate(layout, value string) (time.Time, error) {
	date, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrDateFormat, err.Error())
	}
	return date, nil
}
