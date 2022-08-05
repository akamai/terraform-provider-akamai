package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiApiRequestConstraints_data_basic(t *testing.T) {
	t.Run("match by ApiRequestConstraints ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getAPIRequestConstraintsResponse := appsec.GetApiRequestConstraintsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSApiRequestConstraints/ApiRequestConstraints.json"), &getAPIRequestConstraintsResponse)

		client.On("GetApiRequestConstraints",
			mock.Anything,
			appsec.GetApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1},
		).Return(&getAPIRequestConstraintsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSApiRequestConstraints/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_api_request_constraints.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
