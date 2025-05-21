package networklists

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client networklists.NetworkList, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}
