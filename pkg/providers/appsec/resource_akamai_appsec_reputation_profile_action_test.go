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
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileAction.json")), &cu)

		cr := appsec.GetReputationProfileActionResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResReputationProfileAction/ReputationProfileAction.json")), &cr)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

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
				},
			})
		})

		client.AssertExpectations(t)
	})

}
