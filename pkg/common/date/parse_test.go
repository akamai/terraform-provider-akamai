package date

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		layout        string
		value         string
		expectedError error
	}{
		"ok": {
			layout: DefaultFormat,
			value:  "2016-08-22T23:38:38Z",
		},
		"wrong layout": {
			expectedError: ErrDateFormat,
			layout:        DefaultFormat,
			value:         "2016-22-44T33:88:99Z",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			parsedDate, err := ParseFormat(test.layout, test.value)
			assert.True(t, errors.Is(err, test.expectedError))
			if err == nil {
				assert.Equal(t, test.value, parsedDate.Format(test.layout))
			}
		})
	}
}

func TestFormatRFC3339Nano(t *testing.T) {
	tests := map[string]struct {
		date time.Time
		want string
	}{
		"non zero time": {
			date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "2000-01-01T00:00:00Z",
		},
		"zero time": {
			date: time.Time{},
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := FormatRFC3339Nano(tt.date); got != tt.want {
				t.Errorf("FormatRFC3339Nano() = %v, want %v", got, tt.want)
			}
		})
	}
}
