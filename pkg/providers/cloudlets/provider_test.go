package cloudlets

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type (
	// CustomPollingSubprovider is a cloudlets subprovider with customizable polling intervals for testing.
	CustomPollingSubprovider struct {
		Subprovider
	}
)

// NewCustomPollingSubprovider creates a cloudlets subprovider with custom polling intervals for testing.
func NewCustomPollingSubprovider() *CustomPollingSubprovider {
	return &CustomPollingSubprovider{}
}
