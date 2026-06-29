package clientlists

import (
	"os"
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// testSubproviderConfig returns a subproviderConfig with timing intervals
// that are reduced so tests do not pay the production sleep cost.
func testSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activation: clientListActivationConfig{
			pollActivationInterval:          1 * time.Microsecond,
			activationRetryBaseDelay:        1 * time.Microsecond,
			waitActivationCompletionTimeout: 2 * time.Second,
			activationRetryTimeout:          1 * time.Second,
			activationRetryMaxAttempts:      3,
		},
	}
}

// loadFixtureBytes returns the entire contents of the given file as a byte slice
func loadFixtureBytes(path string) []byte {
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// loadFixtureString returns the entire contents of the given file as a string
func loadFixtureString(path string) string {
	return string(loadFixtureBytes(path))
}
