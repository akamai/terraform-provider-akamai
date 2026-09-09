package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiThreatIntel_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by Threat Intel ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateThreatIntelResponse := appsec.UpdateThreatIntelResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResThreatIntel/ThreatIntel.json"), &updateThreatIntelResponse)
		require.NoError(t, err)

		getThreatIntelResponse := appsec.GetThreatIntelResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResThreatIntel/ThreatIntel.json"), &getThreatIntelResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetThreatIntel",
			testutils.MockContext,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getThreatIntelResponse, nil)

		client.APPSEC.On("UpdateThreatIntel",
			testutils.MockContext,
			appsec.UpdateThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ThreatIntel: "off"},
		).Return(&updateThreatIntelResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThreatIntel/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_threat_intel.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
