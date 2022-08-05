package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiThreatIntel_res_basic(t *testing.T) {
	t.Run("match by Threat Intel ID", func(t *testing.T) {
		client := &mockappsec{}

		updateThreatIntelResponse := appsec.UpdateThreatIntelResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResThreatIntel/ThreatIntel.json"), &updateThreatIntelResponse)

		getThreatIntelResponse := appsec.GetThreatIntelResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResThreatIntel/ThreatIntel.json"), &getThreatIntelResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetThreatIntel",
			mock.Anything,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getThreatIntelResponse, nil)

		client.On("UpdateThreatIntel",
			mock.Anything,
			appsec.UpdateThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ThreatIntel: "off"},
		).Return(&updateThreatIntelResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResThreatIntel/match_by_id.tf"),
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
