package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiApiEndpoints_data_basic(t *testing.T) {
	t.Run("match by ApiEndpoints ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		cv := appsec.GetApiEndpointsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSApiEndpoints/ApiEndpoints.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetApiEndpoints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetApiEndpointsRequest{ConfigID: 43253, Version: 7},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSApiEndpoints/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_api_endpoints.test", "id", "619183"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
