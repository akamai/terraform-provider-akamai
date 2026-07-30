package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimRightWhitespace(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"no trailing whitespace": {
			input:    "hello world",
			expected: "hello world",
		},
		"trailing spaces": {
			input:    "hello world   ",
			expected: "hello world",
		},
		"trailing tabs": {
			input:    "hello world\t\t",
			expected: "hello world",
		},
		"trailing newlines": {
			input:    "hello world\n\n",
			expected: "hello world",
		},
		"mixed trailing whitespace": {
			input:    "hello world \t\n\r ",
			expected: "hello world",
		},
		"only whitespace": {
			input:    " \t\n\r ",
			expected: "",
		},
		"empty string": {
			input:    "",
			expected: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := TrimRightWhitespace(test.input)
			if result != test.expected {
				t.Errorf("expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestImportIDSplitter(t *testing.T) {
	tests := map[string]struct {
		input           string
		formatHint      string
		acceptedLengths []int
		expectedParts   []string
		expectError     string
	}{
		"accepting one part": {
			input:           "12345",
			formatHint:      "resourceID",
			acceptedLengths: []int{1},
			expectedParts:   []string{"12345"},
		},
		"accepting two parts": {
			input:           "12345,grp_234",
			formatHint:      "resourceID,groupID",
			acceptedLengths: []int{2},
			expectedParts:   []string{"12345", "grp_234"},
		},
		"accepting one or two parts - one provided": {
			input:           "12345",
			formatHint:      "resourceID[,groupID]",
			acceptedLengths: []int{1, 2},
			expectedParts:   []string{"12345"},
		},
		"accepting one or two parts - two provided": {
			input:           "12345,678",
			formatHint:      "resourceID[,groupID]",
			acceptedLengths: []int{1, 2},
			expectedParts:   []string{"12345", "678"},
		},
		"valid three parts with spaces": {
			input:           " 555 , 321 , false ",
			formatHint:      "resourceID,groupID,flag",
			acceptedLengths: []int{3},
			expectedParts:   []string{"555", "321", "false"},
		},
		"error - invalid empty input": {
			input:           "",
			formatHint:      "resourceID,groupID",
			acceptedLengths: []int{2},
			expectError:     "importID cannot be empty; you need to provide an importID in the format 'resourceID,groupID'",
		},
		"error - invalid whitespace-only input": {
			input:           "   \t\r\n",
			formatHint:      "resourceID,groupID",
			acceptedLengths: []int{2},
			expectError:     "importID cannot be empty; you need to provide an importID in the format 'resourceID,groupID'",
		},
		"error - accepting two parts but one provided": {
			input:           "123,",
			formatHint:      "resourceID,groupID",
			acceptedLengths: []int{2},
			expectError:     "importID part 2 cannot be empty; you need to provide an importID in the format 'resourceID,groupID'",
		},
		"error - accepting two parts but three provided": {
			input:           "123,456,789",
			formatHint:      "resourceID,groupID",
			acceptedLengths: []int{2},
			expectError:     "invalid number of importID parts: 3; you need to provide an importID in the format 'resourceID,groupID'",
		},
		"error - accepting one or two parts but three provided": {
			input:           "12345,678,true",
			formatHint:      "resourceID[,groupID]",
			acceptedLengths: []int{1, 2},
			expectError:     "invalid number of importID parts: 3; you need to provide an importID in the format 'resourceID[,groupID]'",
		},
		"error - no format hint": {
			input:           "12345,grp_234",
			formatHint:      "",
			acceptedLengths: []int{2},
			expectError:     "no format hint defined for importID; you need to provide a format hint using ImportIDSplitter method",
		},
		"error - no accepted lengths defined": {
			input:           "12345,grp_234",
			formatHint:      "resourceID,groupID",
			acceptedLengths: []int{},
			expectError:     "no accepted lengths defined for importID; you need to provide at least one accepted length using AcceptLen method",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			splitter := ImportIDSplitter(test.formatHint)
			for _, length := range test.acceptedLengths {
				splitter = splitter.AcceptLen(length)
			}
			parts, err := splitter.Split(test.input)
			if test.expectError != "" {
				assert.EqualError(t, err, test.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedParts, parts)
			}
		})
	}
}

func TestToStrings(t *testing.T) {
	t.Parallel()

	type myString string

	tests := map[string]struct {
		input    []myString
		expected []string
	}{
		"multiple elements": {
			input:    []myString{"foo", "bar", "baz"},
			expected: []string{"foo", "bar", "baz"},
		},
		"single element": {
			input:    []myString{"only"},
			expected: []string{"only"},
		},
		"empty slice": {
			input:    []myString{},
			expected: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, ToStrings(tc.input))
		})
	}
}

func TestJoinStringBased(t *testing.T) {
	t.Parallel()

	type myString string

	tests := map[string]struct {
		input    []myString
		sep      string
		expected string
	}{
		"multiple elements with comma separator": {
			input:    []myString{"foo", "bar", "baz"},
			sep:      ", ",
			expected: "foo, bar, baz",
		},
		"single element": {
			input:    []myString{"only"},
			sep:      ", ",
			expected: "only",
		},
		"empty slice": {
			input:    []myString{},
			sep:      ", ",
			expected: "",
		},
		"multiple elements with custom separator": {
			input:    []myString{"a", "b", "c"},
			sep:      "', '",
			expected: "a', 'b', 'c",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, JoinStringBased(tc.input, tc.sep))
		})
	}
}
