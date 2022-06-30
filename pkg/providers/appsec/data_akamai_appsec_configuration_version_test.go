package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfigurationVersion_data_basic(t *testing.T) {
	t.Run("match by ConfigurationVersion ID", func(t *testing.T) {
		client := &mockappsec{}

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSConfigurationVersion/ConfigurationVersion.json"), &getConfigurationVersionsResponse)

		client.On("GetConfigurationVersions",
			mock.Anything,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSConfigurationVersion/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_configuration_version.test", "id", "0"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
