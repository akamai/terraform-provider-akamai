package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalRule_res_basic(t *testing.T) {
	t.Run("match by EvalRule ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateEvalRuleResponse := appsec.UpdateEvalRuleResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalRule/EvalRuleUpdated.json"), &updateEvalRuleResponse)
		require.NoError(t, err)

		getEvalRuleResponse := appsec.GetEvalRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalRule/EvalRule.json"), &getEvalRuleResponse)
		require.NoError(t, err)

		removeEvalRuleActionResponse := appsec.UpdateEvalRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalRule/EvalRule.json"), &removeEvalRuleActionResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalRule",
			testutils.MockContext,
			appsec.GetEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getEvalRuleResponse, nil)

		conditionExceptionJSON := testutils.LoadFixtureBytes(t, "testdata/TestResEvalRule/ConditionException.json")
		client.On("UpdateEvalRule", testutils.MockContext,
			appsec.UpdateEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", RuleID: 12345, JsonPayloadRaw: conditionExceptionJSON},
		).Return(&updateEvalRuleResponse, nil)

		client.On("UpdateEvalRule",
			testutils.MockContext,
			appsec.UpdateEvalRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Action: "none"},
		).Return(&removeEvalRuleActionResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResEvalRule/match_by_id.tf"),
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
