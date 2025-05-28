package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiThreatIntel_res_basic(t *testing.T) {
	t.Run("match by Threat Intel ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateThreatIntelResponse := appsec.UpdateThreatIntelResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResThreatIntel/ThreatIntel.json"), &updateThreatIntelResponse)
		require.NoError(t, err)

		getThreatIntelResponse := appsec.GetThreatIntelResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResThreatIntel/ThreatIntel.json"), &getThreatIntelResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetThreatIntel",
			testutils.MockContext,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getThreatIntelResponse, nil)

		client.On("UpdateThreatIntel",
			testutils.MockContext,
			appsec.UpdateThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ThreatIntel: "off"},
		).Return(&updateThreatIntelResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThreatIntel/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_threat_intel.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
