package slicesets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtract(t *testing.T) {
	var intTests = []struct {
		a, b, want []int
	}{
		{
			a:    []int{},
			b:    []int{},
			want: nil,
		},
		{
			a:    []int{1, 2, 3},
			b:    []int{2, 3},
			want: []int{1},
		},
		{
			a:    nil,
			b:    []int{2, 3},
			want: nil,
		},
		{
			a:    []int{1, 2, 3},
			b:    nil,
			want: []int{1, 2, 3},
		},
	}

	for _, test := range intTests {
		assert.Equal(t, test.want, Subtract(test.a, test.b))
	}

	var stringTests = []struct {
		a, b, want []string
	}{
		{
			a:    []string{},
			b:    []string{},
			want: nil,
		},
		{
			a:    []string{"a", "b", "c"},
			b:    []string{"a", "b"},
			want: []string{"c"},
		},
	}
	for _, test := range stringTests {
		assert.Equal(t, test.want, Subtract(test.a, test.b))
	}
}
