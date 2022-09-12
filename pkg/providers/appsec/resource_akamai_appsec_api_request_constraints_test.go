package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiApiRequestConstraints_res_basic(t *testing.T) {
	t.Run("match by ApiRequestConstraints ID", func(t *testing.T) {
		client := &mockappsec{}

		updateResponse := appsec.UpdateApiRequestConstraintsResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetApiRequestConstraintsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json"), &getResponse)
		require.NoError(t, err)

		deleteResponse := appsec.RemoveApiRequestConstraintsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json"), &deleteResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetApiRequestConstraints",
			mock.Anything,
			appsec.GetApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1},
		).Return(&getResponse, nil)

		client.On("UpdateApiRequestConstraints",
			mock.Anything,
			appsec.UpdateApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1, Action: "alert"},
		).Return(&updateResponse, nil)

		client.On("RemoveApiRequestConstraints",
			mock.Anything,
			appsec.RemoveApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1, Action: "none"},
		).Return(&deleteResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResApiRequestConstraints/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_api_request_constraints.test", "id", "43253:AAAA_81230:1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
