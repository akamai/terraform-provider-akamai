package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomRule_res_basic(t *testing.T) {
	t.Run("CustomRule_basic", func(t *testing.T) {
		client := &appsec.Mock{}

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)
		require.NoError(t, err)

		updateCustomRuleResponse := appsec.UpdateCustomRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &updateCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleAfterUpdate := appsec.GetCustomRuleResponse{} // custom rule after update
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &getCustomRuleAfterUpdate)
		require.NoError(t, err)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)
		require.NoError(t, err)

		getCustomRulesAfterDelete := appsec.GetCustomRulesResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json"), &getCustomRulesAfterDelete)
		require.NoError(t, err)

		client.On("GetCustomRules",
			mock.Anything,
			appsec.GetCustomRulesRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRulesAfterDelete, nil)

		// mock 3 calls to GetCustomRule: 1) after create; 2) via TestCheckResourceAttr 3) pre-update
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponse, nil).Times(3)

		// mock the GetCustomRule call that follows UpdateCustomRule
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleAfterUpdate, nil)

		updateCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/UpdateCustomRule.json")
		client.On("UpdateCustomRule",
			mock.Anything,
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(&updateCustomRuleResponse, nil)

		createCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/CreateCustomRule.json")
		client.On("CreateCustomRule",
			mock.Anything,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&createCustomRuleResponse, nil)

		client.On("RemoveCustomRule",
			mock.Anything,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&removeCustomRuleResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomRule/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCustomRule/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiCustomRule_res_error_removing_active_rule(t *testing.T) {
	t.Run("CustomRule_removing_active_rule", func(t *testing.T) {
		client := &appsec.Mock{}

		updateCustomRuleResponse := appsec.UpdateCustomRuleResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &updateCustomRuleResponse)
		require.NoError(t, err)

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)
		require.NoError(t, err)

		getCustomRuleResponseAfterUpdate := appsec.GetCustomRuleResponse{} // custom rule after update
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &getCustomRuleResponseAfterUpdate)
		require.NoError(t, err)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)
		require.NoError(t, err)

		getCustomRulesAfterDelete := appsec.GetCustomRulesResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json"), &getCustomRulesAfterDelete)
		require.NoError(t, err)

		client.On("GetCustomRules",
			mock.Anything,
			appsec.GetCustomRulesRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRulesAfterDelete, nil)

		// mock 3 calls to GetCustomRule: 1) after create; 2) via TestCheckResourceAttr 3) pre-update
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponse, nil).Times(3)

		// mock the GetCustomRule call that follows UpdateCustomRule
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&getCustomRuleResponseAfterUpdate, nil)

		updateCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/UpdateCustomRule.json")
		client.On("UpdateCustomRule",
			mock.Anything,
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(nil, fmt.Errorf("RemoveCustomRule request failed"))

		createCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/CreateCustomRule.json")
		client.On("CreateCustomRule",
			mock.Anything,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&createCustomRuleResponse, nil)

		client.On("RemoveCustomRule",
			mock.Anything,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&removeCustomRuleResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomRule/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCustomRule/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
						),
						ExpectError: regexp.MustCompile(`RemoveCustomRule request failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
