package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalRuleAction_res_basic(t *testing.T) {
	t.Run("match by EvalRuleAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateEvalRuleActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResEvalRuleAction/EvalRuleActionUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetEvalRuleActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResEvalRuleAction/EvalRuleAction.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetEvalRuleAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 699989},
		).Return(&cr, nil)

		client.On("UpdateEvalRuleAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 699989, Action: "none"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalRuleAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_rule_action.test", "id", "43253:7:AAAA_81230:699989"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
