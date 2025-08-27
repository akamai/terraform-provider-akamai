package iam

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client iam.IAM, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}

// useClient swaps out the client on the global instance for the duration of the given func using both IAM and PAPI
func useIAMandPAPIClient(client iam.IAM, papiClient papi.PAPI, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	origPapi := inst.papiClient
	inst.papiClient = papiClient

	defer func() {
		client = orig
		papiClient = origPapi
		clientLock.Unlock()
	}()

	f()
}
