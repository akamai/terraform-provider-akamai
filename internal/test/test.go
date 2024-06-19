// Package test contains utility code used in tests
package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// NewTimeFromString returns a time value parsed from a string
// in the RFC3339Nano format
func NewTimeFromString(t *testing.T, s string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	require.NoError(t, err)
	return parsedTime
}
