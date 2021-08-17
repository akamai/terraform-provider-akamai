package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationProfileAction_res_basic(t *testing.T) {
	t.Run("match by ReputationProfileAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateReputationProfileActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileAction.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetReputationProfileActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileAction.json"))
		json.Unmarshal([]byte(expectJS), &cr)
		/*
			crd := appsec.UpdateReputationProfileActionResponse{}
			expectJSD := compactJSON(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileActionDelete.json"))
			json.Unmarshal([]byte(expectJSD), &crd)
		*/

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetReputationProfileAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099},
		).Return(&cr, nil)

		client.On("UpdateReputationProfileAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099, Action: "none"},
		).Return(&cu, nil)
		/*
			client.On("RemoveReputationProfileAction",
				mock.Anything, // ctx is irrelevant for this test
				appsec.UpdateReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099, Action: "none"},
			).Return(&crd, nil)
		*/
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationProfileAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_action.test", "id", "43253:AAAA_81230:1685099"),
						),
					},

					/*{
						Config: loadFixtureString("testdata/TestResReputationProfileAction/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_action.test", "id", "43253"),
						),
					},*/
				},
			})
		})

		client.AssertExpectations(t)
	})

}
