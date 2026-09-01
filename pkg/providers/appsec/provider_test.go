package appsec

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// mockGetConfigurationVersionDefault adds a default .Maybe() GetConfigurationVersion mock
// to the client. This replicates the behaviour that useClient() used to provide automatically.
func mockGetConfigurationVersionDefault(client *appsec.Mock) {
	defaultVersionResp := &appsec.GetConfigurationVersionResponse{
		ConfigID:   43253,
		ConfigName: "Akamai Tools",
		Version:    7,
		Production: appsec.EnvironmentStatus{Status: "Inactive"},
		Staging:    appsec.EnvironmentStatus{Status: "Inactive"},
	}
	client.On("GetConfigurationVersion", testutils.MockContext, mock.Anything).Return(defaultVersionResp, nil).Maybe()
}
