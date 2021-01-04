package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRatePolicyAction_res_basic(t *testing.T) {
	t.Run("match by RatePolicyAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateRatePolicyActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRatePolicyAction/RatePolicyUpdateResponse.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetRatePolicyActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRatePolicyAction/RatePolicyActions.json"))
		json.Unmarshal([]byte(expectJS), &cr)
		/*
			cd := appsec.GetRatePolicyActionResponse{}
			expectJSD := compactJSON(loadFixtureBytes("testdata/TestResRatePolicyAction/RatePolicyActionDeleted.json"))
			json.Unmarshal([]byte(expectJSD), &cd)
		*/
		client.On("GetRatePolicyAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRatePolicyActionRequest{ConfigID: 43253, Version: 15, PolicyID: "AAAA_81230", ID: 135355}, //, Ipv4Action: "none", Ipv6Action: "none"},
		).Return(&cr, nil)

		client.On("UpdateRatePolicyAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 15, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "none", Ipv6Action: "none"},
		).Return(&cu, nil)
		/*
			client.On("RemoveRatePolicyAction",
				mock.Anything, // ctx is irrelevant for this test
				appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 15, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "none", Ipv6Action: "alert"},
			).Return(&cd, nil)
		*/
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRatePolicyAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "135355"),
							//resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "none"),
						),
						ExpectNonEmptyPlan: true,
					},

					{
						Config: loadFixtureString("testdata/TestResRatePolicyAction/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "135355"),
							//resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "none"),
						),
						ExpectNonEmptyPlan: true,
					}, /*
						{
							Config: loadFixtureString("testdata/TestResRatePolicyAction/delete_by_id.tf"),
							Check: resource.ComposeTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "321456"),
							//	resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "none"),
							),
							//ExpectNonEmptyPlan: true,
						},*/
				},
			})
		})

		client.AssertExpectations(t)
	})

}
