package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
)

// TODO marks a test as being in a "pending" state and logs a message telling the user why. Such tests are expected to
// fail for the time being and may exist for the sake of unfinished/future features or to document known failure cases
// that won't be fixed right away. The failure of a pending test is not considered an error and the test will therefore
// be skipped unless the TEST_TODO environment variable is set to a non-empty value.
func TODO(t *testing.T, message string) {
	t.Helper()
	t.Log(fmt.Sprintf("TODO: %s (%s)", message, t.Name()))

	if os.Getenv("TEST_TODO") == "" {
		t.Skip("TODO: Set TEST_TODO=1 in env to run this test")
	}
}

// MuteLogging globally prevents logging output unless TEST_LOGGING env var is not empty
func MuteLogging(t *testing.T) {
	t.Helper()

	if os.Getenv("TEST_LOGGING") == "" {
		hclog.SetDefault(hclog.NewNullLogger())
		t.Log("Logging is suppressed. Set TEST_LOGGING=1 in env to see logged messages during test")
	}
}
