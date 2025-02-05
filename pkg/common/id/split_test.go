package id

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	tests := map[string]struct {
		id             string
		expectedLength int
		example        string
		expectedValues []string
		expectErr      error
	}{
		"two parts": {
			id:             "TEST:VALUE",
			example:        "test:value",
			expectedLength: 2,
			expectedValues: []string{"TEST", "VALUE"},
		},
		"three parts with number": {
			id:             "TEST:VALUE:5",
			example:        "test:value:number",
			expectedLength: 3,
			expectedValues: []string{"TEST", "VALUE", "5"},
		},
		"expect error with example": {
			id:             "TEST:VALUE:5:ERROR",
			expectedLength: 3,
			example:        "test:value:number",
			expectErr:      fmt.Errorf("id 'TEST:VALUE:5:ERROR' is incorrectly formatted: should be of form 'test:value:number'"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			parts, err := Split(test.id, test.expectedLength, test.example)
			if test.expectErr != nil {
				assert.EqualError(t, err, test.expectErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedValues, parts)
			}

		})
	}
}
