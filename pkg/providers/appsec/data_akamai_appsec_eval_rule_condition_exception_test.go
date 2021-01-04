package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalRuleConditionException_data_basic(t *testing.T) {
	t.Run("match by EvalRuleConditionException ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetEvalRuleConditionExceptionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSEvalRuleConditionException/EvalRuleConditionException.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetEvalRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSEvalRuleConditionException/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_rule_condition_exception.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
