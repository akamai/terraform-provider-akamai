package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalRuleConditionException_res_basic(t *testing.T) {
	t.Run("match by EvalRuleConditionException ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateEvalRuleConditionExceptionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResEvalRuleConditionException/EvalRuleConditionExceptionUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetEvalRuleConditionExceptionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResEvalRuleConditionException/EvalRuleConditionException.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crr := appsec.RemoveEvalRuleConditionExceptionResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResEvalRuleConditionException/EvalRuleConditionException.json"))
		json.Unmarshal([]byte(expectJSR), &crr)

		client.On("GetEvalRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&cr, nil)

		client.On("RemoveEvalRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveEvalRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&crr, nil)

		client.On("UpdateEvalRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Conditions: []struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{}, Exception: struct {
				HeaderCookieOrParamValues        []string "json:\"headerCookieOrParamValues\""
				SpecificHeaderCookieOrParamNames []struct {
					Names    []string "json:\"names\""
					Selector string   "json:\"selector\""
				} "json:\"specificHeaderCookieOrParamNames\""
			}{HeaderCookieOrParamValues: []string{"abc"}, SpecificHeaderCookieOrParamNames: []struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}(nil)}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalRuleConditionException/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_rule_condition_exception.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
