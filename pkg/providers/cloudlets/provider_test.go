package cloudlets

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// testPollActivationInterval is the polling interval for activation status checks in tests.
	testPollActivationInterval = 1 * time.Millisecond
	// testPollRetryInterval is the polling interval for retrying policy activation in tests.
	testPollRetryInterval = 1 * time.Millisecond
	// testRetryTimeout is the timeout for policy activation retries in tests.
	testRetryTimeout = 1 * time.Millisecond
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type (
	// CustomPollingSubprovider is a cloudlets subprovider with customizable polling intervals for testing.
	CustomPollingSubprovider struct {
		Subprovider
		pollActivationInterval time.Duration
		pollRetryInterval      time.Duration
		retryTimeout           time.Duration
	}
)

// NewCustomPollingSubprovider creates a cloudlets subprovider with custom polling intervals for testing.
func NewCustomPollingSubprovider(pollActivationInterval, pollRetryInterval, retryTimeout time.Duration) *CustomPollingSubprovider {
	return &CustomPollingSubprovider{
		pollActivationInterval: pollActivationInterval,
		pollRetryInterval:      pollRetryInterval,
		retryTimeout:           retryTimeout,
	}
}

// SDKResources overrides the embedded Subprovider's SDKResources to use test polling intervals.
func (p *CustomPollingSubprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cloudlets_application_load_balancer":            resourceCloudletsApplicationLoadBalancer(),
		"akamai_cloudlets_application_load_balancer_activation": resourceCloudletsApplicationLoadBalancerActivation(),
		"akamai_cloudlets_policy":                               resourceCloudletsPolicy(),
		"akamai_cloudlets_policy_activation":                    resourceCloudletsPolicyActivation(p.pollActivationInterval, p.pollRetryInterval, p.retryTimeout),
	}
}
