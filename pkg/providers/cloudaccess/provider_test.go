package cloudaccess

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(cloudAccessClient cloudaccess.CloudAccess, f func()) {
	clientLock.Lock()
	orig := client
	client = cloudAccessClient

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}
