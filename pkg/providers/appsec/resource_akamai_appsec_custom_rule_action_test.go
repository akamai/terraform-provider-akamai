package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomRuleAction_res_basic(t *testing.T) {
	t.Run("match by CustomRuleAction ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &getConfigResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigResponse, nil)

		updateResponse := appsec.UpdateCustomRuleActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRuleAction/CustomRuleActionUpdated.json"), &updateResponse)
		require.NoError(t, err)
		client.On("UpdateCustomRuleAction",
			testutils.MockContext,
			appsec.UpdateCustomRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 60036362, Action: "none"},
		).Return(&updateResponse, nil)

		getResponse := appsec.GetCustomRuleActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRuleAction/CustomRuleAction.json"), &getResponse)
		require.NoError(t, err)
		client.On("GetCustomRuleAction",
			testutils.MockContext,
			appsec.GetCustomRuleActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 60036362},
		).Return(&getResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCustomRuleAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule_action.test", "id", "43253:AAAA_81230:60036362"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
