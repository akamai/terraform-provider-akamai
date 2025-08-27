package cloudlets

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"

	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets/v3"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(cloudletsClient cloudlets.Cloudlets, f func()) {
	clientLock.Lock()
	orig := client
	client = cloudletsClient

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}

// useClientV3 swaps out the client v3 on the global instance for the duration of the given func
func useClientV3(cloudletsV3Client v3.Cloudlets, f func()) {
	clientLock.Lock()
	orig := v3Client
	v3Client = cloudletsV3Client

	defer func() {
		v3Client = orig
		clientLock.Unlock()
	}()

	f()
}

// useClientV2AndV3 swaps out both client (v2) and client v3 on the global instances for the duration of the given func. To be used in by tests for data sources and resources that use both V2 & V3 cloudlets
func useClientV2AndV3(cloudletsV2Client cloudlets.Cloudlets, cloudletsV3Client v3.Cloudlets, f func()) {
	clientLock.Lock()
	origV2 := client
	client = cloudletsV2Client
	origV3 := v3Client
	v3Client = cloudletsV3Client

	defer func() {
		client = origV2
		v3Client = origV3
		clientLock.Unlock()
	}()

	f()
}
