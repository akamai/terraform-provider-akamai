package mtlskeystore

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(mtlsKeystoreClient mtlskeystore.MTLSKeystore, f func()) {
	clientLock.Lock()
	orig := client
	client = mtlsKeystoreClient

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}
