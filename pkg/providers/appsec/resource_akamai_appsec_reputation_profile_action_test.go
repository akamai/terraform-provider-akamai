package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationProfileAction_res_basic(t *testing.T) {
	t.Run("match by ReputationProfileAction ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateReputationProfileActionResponse := appsec.UpdateReputationProfileActionResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileAction.json"), &updateReputationProfileActionResponse)
		require.NoError(t, err)

		getReputationProfileActionResponse := appsec.GetReputationProfileActionResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileAction.json"), &getReputationProfileActionResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetReputationProfileAction",
			mock.Anything,
			appsec.GetReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099},
		).Return(&getReputationProfileActionResponse, nil)

		client.On("UpdateReputationProfileAction",
			mock.Anything,
			appsec.UpdateReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099, Action: "none"},
		).Return(&updateReputationProfileActionResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationProfileAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_action.test", "id", "43253:AAAA_81230:1685099"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
