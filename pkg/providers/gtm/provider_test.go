package gtm

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(gtmClient gtm.GTM, f func()) {
	clientLock.Lock()
	orig := client
	client = gtmClient

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}
