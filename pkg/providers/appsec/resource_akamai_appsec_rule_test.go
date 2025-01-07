package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRule_res_basic(t *testing.T) {
	t.Run("match by Rule ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateRuleResponse := appsec.UpdateRuleResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRule/Rule.json"), &updateRuleResponse)
		require.NoError(t, err)

		getRuleResponse := appsec.GetRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRule/Rule.json"), &getRuleResponse)
		require.NoError(t, err)

		deleteRuleResponse := appsec.UpdateRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRule/Rule.json"), &deleteRuleResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetRule",
			testutils.MockContext,
			appsec.GetRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getRuleResponse, nil)

		conditionExceptionJSON := testutils.LoadFixtureBytes(t, "testdata/TestResRule/ConditionException.json")
		client.On("UpdateRule",
			testutils.MockContext,
			appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", RuleID: 12345, JsonPayloadRaw: conditionExceptionJSON},
		).Return(&updateRuleResponse, nil)

		client.On("UpdateRule",
			testutils.MockContext,
			appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Action: "none"},
		).Return(&deleteRuleResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRule/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rule.test", "id", "43253:AAAA_81230:12345"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
