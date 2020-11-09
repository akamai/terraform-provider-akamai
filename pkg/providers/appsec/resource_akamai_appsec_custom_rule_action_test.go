package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomRuleAction_res_basic(t *testing.T) {
	t.Run("match by CustomRuleAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateCustomRuleActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResCustomRuleAction/CustomRuleActionUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetCustomRuleActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResCustomRuleAction/CustomRuleAction.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetCustomRuleAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 60036362},
		).Return(&cr, nil)

		client.On("UpdateCustomRuleAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateCustomRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 60036362, Action: "none"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomRuleAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule_action.test", "id", "60036362"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
