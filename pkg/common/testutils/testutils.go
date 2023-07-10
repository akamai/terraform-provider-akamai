// Package testutils gathers reusable pieces useful for testing
package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// tfTestTempDir specifies the location of tmp directory which will be used by provider SDK's testing framework
const tfTestTempDir = "./test_tmp"

// TFTestSetup contains common setup for tests in all subproviders
func TFTestSetup() error {
	if err := os.MkdirAll(tfTestTempDir, 0755); err != nil {
		return fmt.Errorf("test setup failed: %s", err)
	}
	if err := os.Setenv("TF_ACC_TEMP_DIR", tfTestTempDir); err != nil {
		return fmt.Errorf("test setup failed: %s", err)
	}
	return nil
}

// TFTestTeardown contains common teardown for tests in all subproviders
func TFTestTeardown() error {
	if err := os.RemoveAll(tfTestTempDir); err != nil {
		return fmt.Errorf("test teardown failed: %s", err)
	}
	if err := os.Unsetenv("TF_ACC_TEMP_DIR"); err != nil {
		return fmt.Errorf("test teardown failed: %s", err)
	}
	return nil
}

// LoadFixtureBytes returns the entire contents of the given file as a byte slice
func LoadFixtureBytes(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return contents
}

// LoadFixtureString returns the entire contents of the given file as a string
func LoadFixtureString(t *testing.T, format string, args ...interface{}) string {
	return string(LoadFixtureBytes(t, fmt.Sprintf(format, args...)))
}
