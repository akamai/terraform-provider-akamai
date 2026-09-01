package dns

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// testSubproviderConfig returns a subproviderConfig with timing intervals
// reduced to milliseconds so tests do not pay the production sleep cost.
func testSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		zone: dnsZoneResourceConfig{
			postCreateChangeListWait:    time.Millisecond,
			preSubmitChangeListWait:     time.Millisecond,
			checkDeletionStatusInterval: time.Millisecond,
		},
		record: dnsRecordResourceConfig{
			conflictRetryInterval:  time.Millisecond,
			soaSerialRetryInterval: time.Millisecond,
		},
	}
}
