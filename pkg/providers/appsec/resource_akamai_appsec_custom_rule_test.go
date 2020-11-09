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

		client.On("GetCustomRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomRuleRequest{ConfigID: 43253, ID: 661699},
		).Return(&cr, nil)

		client.On("UpdateCustomRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateCustomRuleRequest{ConfigID: 43253, ID: 661699, Name: "Rule Test New Updated", Description: "Can I create all conditions?", Version: 0, RuleActivated: false, Tag: []string{"test"}, Conditions: []struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{Type: "requestMethodMatch", PositiveMatch: true, Value: []string{"GET", "CONNECT", "TRACE", "PUT", "POST", "OPTIONS", "DELETE", "HEAD"}, ValueWildcard: false, ValueCase: false, NameWildcard: false, Name: []string(nil), NameCase: false}, struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{Type: "pathMatch", PositiveMatch: true, Value: []string{"/H", "/Li", "/He"}, ValueWildcard: false, ValueCase: false, NameWildcard: false, Name: []string(nil), NameCase: false}, struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{Type: "extensionMatch", PositiveMatch: true, Value: []string{"Li", "He", "H"}, ValueWildcard: true, ValueCase: true, NameWildcard: false, Name: []string(nil), NameCase: false}}},
		).Return(&cu, nil)

		client.On("CreateCustomRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateCustomRuleRequest{ConfigID: 43253, Name: "Rule Test New", Description: "Can I create all conditions?", Version: 0, RuleActivated: false, Tag: []string{"test"}, Conditions: []struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{Type: "requestMethodMatch", PositiveMatch: true, Value: []string{"GET", "CONNECT", "TRACE", "PUT", "POST", "OPTIONS", "DELETE", "HEAD"}, ValueWildcard: false, ValueCase: false, NameWildcard: false, Name: []string(nil), NameCase: false}, struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{Type: "pathMatch", PositiveMatch: true, Value: []string{"/H", "/Li", "/He"}, ValueWildcard: false, ValueCase: false, NameWildcard: false, Name: []string(nil), NameCase: false}, struct {
				Type          string   "json:\"type\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Value         []string "json:\"value,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				ValueCase     bool     "json:\"valueCase,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				Name          []string "json:\"name,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
			}{Type: "extensionMatch", PositiveMatch: true, Value: []string{"Li", "He", "H"}, ValueWildcard: true, ValueCase: true, NameWildcard: false, Name: []string(nil), NameCase: false}}},
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
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "661699"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCustomRule/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "id", "661699"),
							//resource.TestCheckResourceAttr("akamai_appsec_custom_rule.test", "rules", compactJSON(loadFixtureBytes("testdata/TestResCustomRule/custom_rules.json"))),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
