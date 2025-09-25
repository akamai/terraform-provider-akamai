package text

import (
	"fmt"
	"testing"
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
		expectError     error
	}{
		"valid single part": {
			input:           "part1",
			formatHint:      "part1",
			acceptedLengths: []int{1},
			expectedParts:   []string{"part1"},
			expectError:     nil,
		},
		"valid two parts": {
			input:           "part1,part2",
			formatHint:      "part1,part2",
			acceptedLengths: []int{2},
			expectedParts:   []string{"part1", "part2"},
			expectError:     nil,
		},
		"valid three parts with spaces": {
			input:           " part1 , part2 , part3 ",
			formatHint:      "part1,part2,part3",
			acceptedLengths: []int{3},
			expectedParts:   []string{"part1", "part2", "part3"},
			expectError:     nil,
		},
		"invalid empty input": {
			input:           "",
			formatHint:      "part1,part2",
			acceptedLengths: []int{2},
			expectedParts:   nil,
			expectError:     fmt.Errorf("importID cannot be empty; you need to provide an importID in the format 'part1,part2'"),
		},
		"invalid length": {
			input:           "part1,part2,part3",
			formatHint:      "part1,part2",
			acceptedLengths: []int{2},
			expectedParts:   nil,
			expectError:     fmt.Errorf("invalid number of importID parts: '3'; you need to provide an importID in the format 'part1,part2'"),
		},
		"no accepted lengths defined": {
			input:           "part1",
			formatHint:      "part1",
			acceptedLengths: []int{},
			expectedParts:   nil,
			expectError:     fmt.Errorf("no accepted lengths defined for importID; you need to provide at least one accepted length using AcceptLen method"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			splitter := ImportIDSplitter(test.formatHint)
			for _, length := range test.acceptedLengths {
				splitter = splitter.AcceptLen(length)
			}
			parts, err := splitter.Split(test.input)
			if test.expectError != nil {
				if err == nil || err.Error() != test.expectError.Error() {
					t.Errorf("expected error '%v', got '%v'", test.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if len(parts) != len(test.expectedParts) {
					t.Errorf("expected parts '%v', got '%v'", test.expectedParts, parts)
				} else {
					for i := range parts {
						if parts[i] != test.expectedParts[i] {
							t.Errorf("expected parts '%v', got '%v'", test.expectedParts, parts)
							break
						}
					}
				}
			}
		})
	}
}
