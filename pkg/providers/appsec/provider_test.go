package appsec

import (
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client appsec.APPSEC, f func()) {
	clientLock.Lock()
	orig := inst.client

	// If client is a mock, add default GetConfigurationVersion response for the new method
	if mockClient, ok := client.(*appsec.Mock); ok {
		// Create a simple response that works for most tests
		defaultVersionResp := &appsec.GetConfigurationVersionResponse{
			ConfigID:   43253,
			ConfigName: "Akamai Tools",
			Version:    7,
			Production: appsec.EnvironmentStatus{Status: "Inactive"},
			Staging:    appsec.EnvironmentStatus{Status: "Inactive"},
		}
		mockClient.On("GetConfigurationVersion", mock.Anything, mock.Anything).Return(defaultVersionResp, nil).Maybe()
	}

	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}
