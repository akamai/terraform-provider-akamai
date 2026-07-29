package networklists

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// testSubproviderConfig returns a subproviderConfig with timing intervals
// reduced to milliseconds so tests do not pay the production sleep cost.
func testSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activations: resourceActivationsConfig{
			activationPollInterval: 1 * time.Millisecond,
			createActivationRetry:  1 * time.Millisecond,
		},
	}
}
