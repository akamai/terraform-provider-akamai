package mtlstruststore

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient().
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func.
func useClient(mtlsTruststoreClient mtlstruststore.MTLSTruststore, f func()) {
	clientLock.Lock()
	orig := client
	client = mtlsTruststoreClient

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}
