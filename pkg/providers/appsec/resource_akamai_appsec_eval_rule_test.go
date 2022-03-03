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

		cu := appsec.UpdateEvalRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResEvalRule/EvalRuleUpdated.json")), &cu)

		cr := appsec.GetEvalRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResEvalRule/EvalRule.json")), &cr)

		cd := appsec.UpdateEvalRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResEvalRule/EvalRule.json")), &cd)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&cr, nil)

		conditionExceptionJSON := loadFixtureBytes("testdata/TestResEvalRule/ConditionException.json")
		client.On("UpdateEvalRule", mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", RuleID: 12345, JsonPayloadRaw: conditionExceptionJSON},
		).Return(&cu, nil)

		client.On("UpdateEvalRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Action: "none"},
		).Return(&cd, nil)

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
