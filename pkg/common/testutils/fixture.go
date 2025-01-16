package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadFixtureBytes returns the entire contents of the given file as a byte slice
func LoadFixtureBytes(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	return contents
}

// LoadFixtureStringf returns the entire contents of the given file as a string
func LoadFixtureStringf(t *testing.T, format string, args ...interface{}) string {
	return string(LoadFixtureBytes(t, fmt.Sprintf(format, args...)))
}

// LoadFixtureString returns the entire contents of the given file as a string
func LoadFixtureString(t *testing.T, path string) string {
	return string(LoadFixtureBytes(t, path))
}
