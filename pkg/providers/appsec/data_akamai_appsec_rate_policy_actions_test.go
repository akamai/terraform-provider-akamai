package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRatePolicyActions_data_basic(t *testing.T) {
	t.Run("match by RatePolicyActions ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRatePolicyActionsResponse := appsec.GetRatePolicyActionsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSRatePolicyActions/RatePolicyActions.json"), &getRatePolicyActionsResponse)

		client.On("GetRatePolicyActions",
			mock.Anything,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getRatePolicyActionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRatePolicyActions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_rate_policy_actions.test", "id", "102720"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
