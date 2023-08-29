package collections

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEachInSlice(t *testing.T) {
	var intTests = []struct {
		s, want []int
		fn      func(d int) int
	}{
		{
			s:    []int{},
			want: []int{},
			fn:   func(d int) int { return d },
		},
		{
			s:    []int{1, 2, 3},
			want: []int{2, 3, 4},
			fn:   func(d int) int { return d + 1 },
		},
		{
			s:    []int{1, 2, 3},
			want: []int{1, 4, 9},
			fn:   func(d int) int { return d * d },
		},
	}

	for _, test := range intTests {
		ForEachInSlice(test.s, test.fn)
		assert.Equal(t, test.want, test.s)
	}

	var floatTests = []struct {
		s, want []float64
		fn      func(d float64) float64
	}{
		{
			s:    []float64{1, 2, 3},
			want: []float64{1, 2, 3},
			fn:   func(f float64) float64 { return f },
		},
	}

	for _, test := range floatTests {
		ForEachInSlice(test.s, test.fn)
		assert.Equal(t, test.want, test.s)
	}

	var stringTests = []struct {
		s, want []string
		fn      func(s string) string
	}{
		{
			s:    []string{"a", "b", "c"},
			want: []string{"ab", "bb", "cb"},
			fn:   func(s string) string { return s + "b" },
		},
		{
			s:    []string{"prefa", "prefb", "c"},
			want: []string{"a", "b", "c"},
			fn:   func(s string) string { return strings.TrimPrefix(s, "pref") },
		},
	}

	for _, test := range stringTests {
		ForEachInSlice(test.s, test.fn)
		assert.Equal(t, test.want, test.s)
	}

	type mySlice []int
	mySliceTests := []struct {
		s, want mySlice
		fn      func(d int) int
	}{
		{
			s:    mySlice{1, 2, 3},
			want: mySlice{2, 3, 4},
			fn:   func(d int) int { return d + 1 },
		},
		{
			s:    mySlice{1, 2, 3},
			want: mySlice{1, 4, 9},
			fn:   func(d int) int { return d * d },
		},
	}

	for _, test := range mySliceTests {
		ForEachInSlice(test.s, test.fn)
		assert.Equal(t, test.want, test.s)
	}
}
