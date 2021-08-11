package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomRule_res_basic(t *testing.T) {
	t.Run("match by CustomRule ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateCustomRuleResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResCustomRule/CustomRuleUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cc := appsec.CreateCustomRuleResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"))
		json.Unmarshal([]byte(expectJSC), &cc)

		cr := appsec.GetCustomRuleResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResCustomRule/CustomRule.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crr := appsec.RemoveCustomRuleResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesDeleted.json"))
		json.Unmarshal([]byte(expectJSR), &crr)

		cv := appsec.GetCustomRulesResponse{}
		expectJR := compactJSON(loadFixtureBytes("testdata/TestResCustomRule/CustomRulesForDelete.json"))
		json.Unmarshal([]byte(expectJR), &cv)

		client.On("GetCustomRules",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomRulesRequest{ConfigID: 43253, ID: 661699},
		).Return(&cv, nil)

		client.On("GetCustomRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&cr, nil)

		updateCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/UpdateCustomRule.json")
		client.On("UpdateCustomRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Version: 0, JsonPayloadRaw: updateCustomRuleJSON},
		).Return(&cu, nil)

		createCustomRuleJSON := loadFixtureBytes("testdata/TestResCustomRule/CreateCustomRule.json")
		client.On("CreateCustomRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Version: 0, JsonPayloadRaw: createCustomRuleJSON},
		).Return(&cc, nil)

		client.On("RemoveCustomRule",
			mock.Anything, // ctx is irrelevant for this test
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
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResCustomRule/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "43253:661699"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
