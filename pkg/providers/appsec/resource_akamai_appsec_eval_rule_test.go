package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalRule_res_basic(t *testing.T) {
	t.Run("match by EvalRule ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateEvalRuleResponse := appsec.UpdateEvalRuleResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResEvalRule/EvalRuleUpdated.json"), &updateEvalRuleResponse)
		require.NoError(t, err)

		getEvalRuleResponse := appsec.GetEvalRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEvalRule/EvalRule.json"), &getEvalRuleResponse)
		require.NoError(t, err)

		removeEvalRuleActionResponse := appsec.UpdateEvalRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEvalRule/EvalRule.json"), &removeEvalRuleActionResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

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
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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
