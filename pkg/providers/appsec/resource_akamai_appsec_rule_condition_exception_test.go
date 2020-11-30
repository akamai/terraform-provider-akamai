package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRuleConditionException_res_basic(t *testing.T) {
	t.Run("match by RuleConditionException ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateRuleConditionExceptionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRuleConditionException/RuleConditionException.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cur := appsec.RemoveRuleConditionExceptionResponse{}
		expectJSUR := compactJSON(loadFixtureBytes("testdata/TestResRuleConditionException/RuleConditionException.json"))
		json.Unmarshal([]byte(expectJSUR), &cur)

		cr := appsec.GetRuleConditionExceptionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRuleConditionException/RuleConditionException.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&cr, nil)

		client.On("UpdateRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Conditions: []struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "extensionMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string(nil)}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "filenameMatch", Filenames: []string{"test2"}, PositiveMatch: true, Methods: []string(nil)}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "hostMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string(nil)}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "ipMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string(nil)}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "uriQueryMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string(nil)}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "requestHeaderMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string(nil)}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "requestMethodMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string{"GET"}}, struct {
				Type          string   "json:\"type\""
				Filenames     []string "json:\"filenames,omitempty\""
				PositiveMatch bool     "json:\"positiveMatch\""
				Methods       []string "json:\"methods,omitempty\""
			}{Type: "pathMatch", Filenames: []string(nil), PositiveMatch: true, Methods: []string(nil)}}, Exception: struct {
				HeaderCookieOrParamValues        []string "json:\"headerCookieOrParamValues\""
				SpecificHeaderCookieOrParamNames []struct {
					Names    []string "json:\"names\""
					Selector string   "json:\"selector\""
				} "json:\"specificHeaderCookieOrParamNames\""
			}{HeaderCookieOrParamValues: []string{"test"}, SpecificHeaderCookieOrParamNames: []struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}{struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}{Names: []string{"test"}, Selector: "REQUEST_HEADERS"}, struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}{Names: []string{"test"}, Selector: "REQUEST_COOKIES"}, struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}{Names: []string{"test"}, Selector: "ARGS"}, struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}{Names: []string{"test"}, Selector: "JSON_PAIRS"}, struct {
				Names    []string "json:\"names\""
				Selector string   "json:\"selector\""
			}{Names: []string{"test"}, Selector: "XML_PAIRS"}}}},
		).Return(&cu, nil)

		client.On("RemoveRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveRuleConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&cur, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRuleConditionException/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rule_condition_exception.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
