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
	t.Run("match by CustomRule ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json")), &cu)

		cc := appsec.CreateCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json")), &cc)

		cr := appsec.GetCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json")), &cr)

		crau := appsec.GetCustomRuleResponse{} // custom rule after update
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json")), &crau)

		crr := appsec.RemoveCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json")), &crr)

		cv := appsec.GetCustomRulesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json")), &cv)

		client.On("GetCustomRules",
			mock.Anything,
			appsec.GetCustomRulesRequest{ConfigID: 43253, ID: 661699},
		).Return(&cv, nil)

		// mock 3 calls to GetCustomRule: 1) after create; 2) via TestCheckResourceAttr 3) pre-update
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&cr, nil).Times(3)

		// mock the GetCustomRule call that follows UpdateCustomRule
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&crau, nil)

		updateCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/UpdateCustomRule.json")
		client.On("UpdateCustomRule",
			mock.Anything,
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(&cu, nil)

		createCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/CreateCustomRule.json")
		client.On("CreateCustomRule",
			mock.Anything,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&cc, nil)

		client.On("RemoveCustomRule",
			mock.Anything,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&crr, nil)

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
	t.Run("match by CustomRule ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json")), &cu)

		cc := appsec.CreateCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json")), &cc)

		cr := appsec.GetCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json")), &cr)

		crau := appsec.GetCustomRuleResponse{} // custom rule after update
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json")), &crau)

		crr := appsec.RemoveCustomRuleResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json")), &crr)

		cv := appsec.GetCustomRulesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json")), &cv)

		client.On("GetCustomRules",
			mock.Anything,
			appsec.GetCustomRulesRequest{ConfigID: 43253, ID: 661699},
		).Return(&cv, nil)

		// mock 3 calls to GetCustomRule: 1) after create; 2) via TestCheckResourceAttr 3) pre-update
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&cr, nil).Times(3)

		// mock the GetCustomRule call that follows UpdateCustomRule
		client.On("GetCustomRule",
			mock.Anything,
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&crau, nil)

		updateCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/UpdateCustomRule.json")
		client.On("UpdateCustomRule",
			mock.Anything,
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(nil, fmt.Errorf("RemoveCustomRule request failed"))

		createCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/CreateCustomRule.json")
		client.On("CreateCustomRule",
			mock.Anything,
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&cc, nil)

		client.On("RemoveCustomRule",
			mock.Anything,
			appsec.RemoveCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&crr, nil)

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
