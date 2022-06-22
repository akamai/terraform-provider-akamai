package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfiguration_data_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &mockappsec{}

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)

		client.On("GetConfigurations",
			mock.Anything,
			appsec.GetConfigurationsRequest{},
		).Return(&getConfigurationsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
