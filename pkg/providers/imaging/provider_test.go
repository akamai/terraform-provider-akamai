package imaging

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	PolicyDepth = 4
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client imaging.Imaging, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}
