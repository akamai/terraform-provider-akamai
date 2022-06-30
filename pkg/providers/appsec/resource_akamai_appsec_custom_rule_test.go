package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomRule_res_basic(t *testing.T) {
	t.Run("CustomRule_basic", func(t *testing.T) {
		client := &mockappsec{}

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)

		updateCustomRuleResponse := appsec.UpdateCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &updateCustomRuleResponse)

		getCustomRuleAfterUpdate := appsec.GetCustomRuleResponse{} // custom rule after update
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &getCustomRuleAfterUpdate)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)

		getCustomRulesAfterDelete := appsec.GetCustomRulesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json"), &getCustomRulesAfterDelete)

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

func TestAccAkamaiCustomRule_res_error_removing_active_rule(t *testing.T) {
	t.Run("CustomRule_removing_active_rule", func(t *testing.T) {
		client := &mockappsec{}

		updateCustomRuleResponse := appsec.UpdateCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &updateCustomRuleResponse)

		createCustomRuleResponse := appsec.CreateCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &createCustomRuleResponse)

		getCustomRuleResponse := appsec.GetCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"), &getCustomRuleResponse)

		getCustomRuleResponseAfterUpdate := appsec.GetCustomRuleResponse{} // custom rule after update
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"), &getCustomRuleResponseAfterUpdate)

		removeCustomRuleResponse := appsec.RemoveCustomRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json"), &removeCustomRuleResponse)

		getCustomRulesAfterDelete := appsec.GetCustomRulesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json"), &getCustomRulesAfterDelete)

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
