package date

import (
	"errors"
	"testing"

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
