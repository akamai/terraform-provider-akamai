package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiApiHostnameCoverageMatchTargets_data_basic(t *testing.T) {
	t.Run("match by ApiHostnameCoverageMatchTargets ID", func(t *testing.T) {
		client := &mockappsec{}

		getMatchTargetsResponse := appsec.GetApiHostnameCoverageMatchTargetsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSApiHostnameCoverageMatchTargets/ApiHostnameCoverageMatchTargets.json"), &getMatchTargetsResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetApiHostnameCoverageMatchTargets",
			mock.Anything,
			appsec.GetApiHostnameCoverageMatchTargetsRequest{ConfigID: 43253, Version: 7, Hostname: "rinaldi.sandbox.akamaideveloper.com"},
		).Return(&getMatchTargetsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSApiHostnameCoverageMatchTargets/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_hostname_coverage_match_targets.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
