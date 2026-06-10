package gtm

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

const (
	testGTMSleepInterval   = 1 * time.Millisecond
	testGTMDefaultInterval = 1 * time.Millisecond
	testGTMRetryInterval   = 1 * time.Millisecond
	testGTMMaxRetryTimeout = 10 * time.Millisecond
)

func testSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		domain: gtmDomainResourceConfig{
			sleepInterval:   testGTMSleepInterval,
			defaultInterval: testGTMDefaultInterval,
		},
		property: gtmPropertyResourceConfig{
			retryInterval:   testGTMRetryInterval,
			maxRetryTimeout: testGTMMaxRetryTimeout,
			defaultInterval: testGTMDefaultInterval,
		},
		defaultInterval: testGTMDefaultInterval,
	}
}
