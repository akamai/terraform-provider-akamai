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
	// testALBPollActivationInterval is the polling interval for ALB activation status checks in tests.
	testALBPollActivationInterval = 1 * time.Millisecond
	// testALBRetryTimeout is the timeout for ALB activation retries in tests.
	// This needs to be larger than testALBPollActivationInterval to allow retries.
	testALBRetryTimeout = 10 * time.Millisecond
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type (
	// CustomPollingSubprovider is a cloudlets subprovider with customizable polling intervals for testing.
	CustomPollingSubprovider struct {
		Subprovider
		pollActivationInterval    time.Duration
		pollRetryInterval         time.Duration
		retryTimeout              time.Duration
		albPollActivationInterval time.Duration
		albRetryTimeout           time.Duration
	}
)

// NewCustomPollingSubprovider creates a cloudlets subprovider with default test polling intervals.
// Use WithPolicyActivationIntervals and WithALBActivationIntervals to customize specific intervals.
func NewCustomPollingSubprovider() *CustomPollingSubprovider {
	return &CustomPollingSubprovider{
		pollActivationInterval:    defaultPollActivationInterval,
		pollRetryInterval:         defaultPollRetryInterval,
		retryTimeout:              defaultRetryTimeout,
		albPollActivationInterval: defaultALBPollActivationInterval,
		albRetryTimeout:           defaultALBRetryTimeout,
	}
}

// WithPolicyActivationIntervals sets custom polling intervals for policy activation.
func (p *CustomPollingSubprovider) WithPolicyActivationIntervals(pollActivationInterval, pollRetryInterval, retryTimeout time.Duration) *CustomPollingSubprovider {
	p.pollActivationInterval = pollActivationInterval
	p.pollRetryInterval = pollRetryInterval
	p.retryTimeout = retryTimeout
	return p
}

// WithALBActivationIntervals sets custom polling intervals for ALB activation.
func (p *CustomPollingSubprovider) WithALBActivationIntervals(pollActivationInterval, retryTimeout time.Duration) *CustomPollingSubprovider {
	p.albPollActivationInterval = pollActivationInterval
	p.albRetryTimeout = retryTimeout
	return p
}

// SDKResources overrides the embedded Subprovider's SDKResources to use test polling intervals.
func (p *CustomPollingSubprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cloudlets_application_load_balancer":            resourceCloudletsApplicationLoadBalancer(),
		"akamai_cloudlets_application_load_balancer_activation": resourceCloudletsApplicationLoadBalancerActivation(p.albPollActivationInterval, p.albRetryTimeout),
		"akamai_cloudlets_policy":                               resourceCloudletsPolicy(),
		"akamai_cloudlets_policy_activation":                    resourceCloudletsPolicyActivation(p.pollActivationInterval, p.pollRetryInterval, p.retryTimeout),
	}
}
