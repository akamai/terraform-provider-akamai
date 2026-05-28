package accountprotection

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

// mockGetConfigVersion sets up the APPSEC mock to return version 15 for config 43253,
// which is needed because account protection resources call getLatestConfigVersion/getModifiableConfigVersion.
func mockGetConfigVersion(client *edgegrid.TestClient) {
	client.APPSEC.On("GetConfiguration", testutils.MockContext,
		appsec.GetConfigurationRequest{ConfigID: 43253},
	).Return(&appsec.GetConfigurationResponse{
		ID:            43253,
		LatestVersion: 15,
	}, nil).Maybe()
	client.APPSEC.On("GetConfigurationVersion", testutils.MockContext,
		appsec.GetConfigurationVersionRequest{ConfigID: 43253, Version: 15},
	).Return(&appsec.GetConfigurationVersionResponse{
		ConfigID:   43253,
		Version:    15,
		Staging:    appsec.EnvironmentStatus{Status: "Inactive"},
		Production: appsec.EnvironmentStatus{Status: "Inactive"},
	}, nil).Maybe()
}

func compactJSON(message string) string {
	var dst bytes.Buffer
	err := json.Compact(&dst, []byte(message))
	if err != nil {
		panic(err)
	}
	return dst.String()
}
