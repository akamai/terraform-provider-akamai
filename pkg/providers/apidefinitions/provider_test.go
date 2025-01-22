package apidefinitions

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(mockClient apidefinitions.APIDefinitions, mockClientV0 v0.APIDefinitions, f func()) {
	clientLock.Lock()
	orig := client
	client = mockClient
	clientV0 = mockClientV0

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	if f != nil {
		f()
	}
}

var providerConfig = `provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}
`
