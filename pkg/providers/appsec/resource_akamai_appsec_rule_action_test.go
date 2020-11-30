package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRuleAction_res_basic(t *testing.T) {
	t.Run("match by RuleAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateRuleActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRuleAction/RuleAction.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetRuleActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRuleAction/RuleAction.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetRuleAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 699989},
		).Return(&cr, nil)

		client.On("UpdateRuleAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 699989, Action: "none"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRuleAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rule_action.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
