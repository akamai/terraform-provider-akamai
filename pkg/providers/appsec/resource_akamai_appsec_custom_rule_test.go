package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomRule_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("CustomRule_basic", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)
		require.NoError(t, err)

		updateCustomRuleResponse := appsec.UpdateCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRuleUpdated.json"), &updateCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleAfterUpdate := appsec.GetCustomRuleResponse{} // custom rule after update
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRuleUpdated.json"), &getCustomRuleAfterUpdate)
		require.NoError(t, err)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)
		require.NoError(t, err)

		getCustomRulesAfterDelete := appsec.GetCustomRulesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRulesForDelete.json"), &getCustomRulesAfterDelete)
		require.NoError(t, err)

		getCustomRuleUsageRequest := appsec.GetCustomRulesUsageRequest{
			ConfigID: 43253,
			Version:  2,
			RequestBody: appsec.RuleIDs{
				IDs: []int64{661699},
			},
		}

		// mock 3 calls to GetCustomRule: 1) after create; 2) via TestCheckResourceAttr 3) pre-update
		client.APPSEC.On("GetCustomRule",
			testutils.MockContext,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponse, nil).Times(3)

		// mock the GetCustomRule call that follows UpdateCustomRule
		client.APPSEC.On("GetCustomRule",
			testutils.MockContext,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleAfterUpdate, nil)

		updateCustomRuleJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/UpdateCustomRule.json")
		client.APPSEC.On("UpdateCustomRule",
			testutils.MockContext,
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(&updateCustomRuleResponse, nil)

		createCustomRuleJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CreateCustomRule.json")
		client.APPSEC.On("CreateCustomRule",
			testutils.MockContext,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&createCustomRuleResponse, nil)

		client.APPSEC.On("RemoveCustomRule",
			testutils.MockContext,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&removeCustomRuleResponse, nil)

		mockGetLatestConfiguration(client.APPSEC, 43253, 1)
		mockGetCustomRulesUsage(client.APPSEC, getCustomRuleUsageRequest, appsec.GetCustomRulesUsageResponse{Rules: []appsec.CustomRuleUsage{}}, 1)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomRule/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomRule/update_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}

func TestAkamaiCustomRule_res_error_removing_active_rule(t *testing.T) {
	t.Parallel()
	t.Run("CustomRule_removing_active_rule", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateCustomRuleResponse := appsec.UpdateCustomRuleResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRuleUpdated.json"), &updateCustomRuleResponse)
		require.NoError(t, err)

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponseAfterUpdate := appsec.GetCustomRuleResponse{} // custom rule after update
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRuleUpdated.json"), &getCustomRuleResponseAfterUpdate)
		require.NoError(t, err)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)
		require.NoError(t, err)

		getCustomRulesAfterDelete := appsec.GetCustomRulesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRulesForDelete.json"), &getCustomRulesAfterDelete)
		require.NoError(t, err)

		getCustomRuleUsageRequest := appsec.GetCustomRulesUsageRequest{
			ConfigID: 43253,
			Version:  2,
			RequestBody: appsec.RuleIDs{
				IDs: []int64{661699},
			},
		}

		// mock 3 calls to GetCustomRule: 1) after create; 2) via TestCheckResourceAttr 3) pre-update
		client.APPSEC.On("GetCustomRule",
			testutils.MockContext,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponse, nil).Times(3)

		// mock the GetCustomRule call that follows UpdateCustomRule
		client.APPSEC.On("GetCustomRule",
			testutils.MockContext,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponseAfterUpdate, nil)

		updateCustomRuleJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/UpdateCustomRule.json")
		client.APPSEC.On("UpdateCustomRule",
			testutils.MockContext,
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(nil, fmt.Errorf("RemoveCustomRule request failed"))

		createCustomRuleJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CreateCustomRule.json")
		client.APPSEC.On("CreateCustomRule",
			testutils.MockContext,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&createCustomRuleResponse, nil)

		client.APPSEC.On("RemoveCustomRule",
			testutils.MockContext,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&removeCustomRuleResponse, nil)

		mockGetLatestConfiguration(client.APPSEC, 43253, 1)
		mockGetCustomRulesUsage(client.APPSEC, getCustomRuleUsageRequest, appsec.GetCustomRulesUsageResponse{Rules: []appsec.CustomRuleUsage{}}, 1)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomRule/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomRule/update_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
					),
					ExpectError: regexp.MustCompile(`RemoveCustomRule request failed`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}

func TestAkamaiCustomRule_res_error_deleting_rule_in_use(t *testing.T) {
	t.Parallel()
	t.Run("CustomRule_removing_rule_in_use", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)
		require.NoError(t, err)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetCustomRule",
			testutils.MockContext,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponse, nil).Times(3)

		createCustomRuleJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomRule/CreateCustomRule.json")
		client.APPSEC.On("CreateCustomRule",
			testutils.MockContext,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&createCustomRuleResponse, nil)

		mockGetLatestConfiguration(client.APPSEC, 43253, 2)
		client.APPSEC.On("GetCustomRulesUsage",
			testutils.MockContext,
			appsec.GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  2,
				RequestBody: appsec.RuleIDs{
					IDs: []int64{661699},
				},
			},
		).Return(&appsec.GetCustomRulesUsageResponse{
			Rules: []appsec.CustomRuleUsage{
				{
					RuleID: 661699,
					Policies: []appsec.Policy{
						{PolicyID: "p1", PolicyName: "Policy One"},
						{PolicyID: "p2", PolicyName: "Policy Two"},
					},
				},
			},
		}, nil).Once()

		client.APPSEC.On("GetCustomRulesUsage",
			testutils.MockContext,
			appsec.GetCustomRulesUsageRequest{
				ConfigID: 43253,
				Version:  2,
				RequestBody: appsec.RuleIDs{
					IDs: []int64{661699},
				},
			},
		).Return(&appsec.GetCustomRulesUsageResponse{
			Rules: []appsec.CustomRuleUsage{},
		}, nil).Once()

		client.APPSEC.On("RemoveCustomRule",
			testutils.MockContext,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&removeCustomRuleResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomRule/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
					),
				},
				testutils.TestStepDestroyFailed(testutils.LoadFixtureString(t, "testdata/TestResCustomRule/match_by_id.tf"),
					regexp.MustCompile(`custom rule with ID: 661699 cannot be deleted, it is either active or in use in the security policies with IDs: p1, p2`)),
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}
