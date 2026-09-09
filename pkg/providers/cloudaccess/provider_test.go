package cloudaccess

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	testPollingInterval   = 1 * time.Millisecond
	testActivationTimeout = 20 * time.Millisecond
	testUpdateTimeout     = 20 * time.Millisecond
	testDeleteTimeout     = 40 * time.Millisecond
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// testSubprovider is a Subprovider variant with configurable timing for KeyResource.
type testSubprovider struct {
	Subprovider
	activationTimeout time.Duration
	updateTimeout     time.Duration
	deleteTimeout     time.Duration
	pollingInterval   time.Duration
}

// newTestSubprovider creates a testSubprovider with default test-friendly timing.
func newTestSubprovider() *testSubprovider {
	return &testSubprovider{
		activationTimeout: testActivationTimeout,
		updateTimeout:     testUpdateTimeout,
		deleteTimeout:     testDeleteTimeout,
		pollingInterval:   testPollingInterval,
	}
}

// withTimeouts returns the testSubprovider with overridden timing values.
func (p *testSubprovider) withTimeouts(activationTimeout, updateTimeout, deleteTimeout, pollingInterval time.Duration) *testSubprovider {
	p.activationTimeout = activationTimeout
	p.updateTimeout = updateTimeout
	p.deleteTimeout = deleteTimeout
	p.pollingInterval = pollingInterval
	return p
}

// FrameworkResources overrides the embedded Subprovider to inject test-friendly timing into KeyResource.
func (p *testSubprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &KeyResource{
				activationTimeout: p.activationTimeout,
				updateTimeout:     p.updateTimeout,
				deleteTimeout:     p.deleteTimeout,
				pollingInterval:   p.pollingInterval,
			}
		},
	}
}
