package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalRule_res_basic(t *testing.T) {
	t.Run("match by EvalRule ID", func(t *testing.T) {
		client := &mockappsec{}

		updateEvalRuleResponse := appsec.UpdateEvalRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResEvalRule/EvalRuleUpdated.json"), &updateEvalRuleResponse)

		getEvalRuleResponse := appsec.GetEvalRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResEvalRule/EvalRule.json"), &getEvalRuleResponse)

		removeEvalRuleActionResponse := appsec.UpdateEvalRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResEvalRule/EvalRule.json"), &removeEvalRuleActionResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalRule",
			mock.Anything,
			appsec.GetEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getEvalRuleResponse, nil)

		conditionExceptionJSON := loadFixtureBytes("testdata/TestResEvalRule/ConditionException.json")
		client.On("UpdateEvalRule", mock.Anything,
			appsec.UpdateEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", RuleID: 12345, JsonPayloadRaw: conditionExceptionJSON},
		).Return(&updateEvalRuleResponse, nil)

		client.On("UpdateEvalRule",
			mock.Anything,
			appsec.UpdateEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Action: "none"},
		).Return(&removeEvalRuleActionResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalRule/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_rule.test", "id", "43253:AAAA_81230:12345"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
