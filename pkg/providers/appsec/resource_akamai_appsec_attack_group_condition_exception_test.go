package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroupConditionException_res_basic(t *testing.T) {
	t.Run("match by AttackGroupConditionException ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAttackGroupConditionExceptionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAttackGroupConditionException/AttackGroupConditionException.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAttackGroupConditionExceptionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAttackGroupConditionException/AttackGroupConditionException.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crr := appsec.RemoveAttackGroupConditionExceptionResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResAttackGroupConditionException/AttackGroupConditionException.json"))
		json.Unmarshal([]byte(expectJSR), &crr)

		client.On("GetAttackGroupConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&cr, nil)

		client.On("RemoveAttackGroupConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveAttackGroupConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&crr, nil)

		client.On("UpdateAttackGroupConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Conditions: []struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{}, Exception: struct {
				HeaderCookieOrParamValues        []string "json:\"headerCookieOrParamValues\""
				SpecificHeaderCookieOrParamNames []struct {
					Names    []string "json:\"names,omitempty\""
					Selector string   "json:\"selector,omitempty\""
				} "json:\"specificHeaderCookieOrParamNames,omitempty\""
				SpecificHeaderCookieOrParamPrefix struct {
					Prefix   string "json:\"prefix,omitempty\""
					Selector string "json:\"selector,omitempty\""
				} "json:\"specificHeaderCookieOrParamPrefix,omitempty\""
			}{HeaderCookieOrParamValues: []string{"abc"}, SpecificHeaderCookieOrParamNames: []struct {
				Names    []string "json:\"names,omitempty\""
				Selector string   "json:\"selector,omitempty\""
			}(nil), SpecificHeaderCookieOrParamPrefix: struct {
				Prefix   string "json:\"prefix,omitempty\""
				Selector string "json:\"selector,omitempty\""
			}{Prefix: "a*", Selector: "REQUEST_COOKIES"}}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAttackGroupConditionException/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_attack_group_condition_exception.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
