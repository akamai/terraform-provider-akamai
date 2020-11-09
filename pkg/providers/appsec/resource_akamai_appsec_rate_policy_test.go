package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRatePolicy_res_basic(t *testing.T) {
	t.Run("match by RatePolicy ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateRatePolicyResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetRatePolicyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicy.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crpol := appsec.CreateRatePolicyResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicy.json"))
		json.Unmarshal([]byte(expectJSC), &crpol)

		crpolr := appsec.RemoveRatePolicyResponse{}
		expectJSCR := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyEmpty.json"))
		json.Unmarshal([]byte(expectJSCR), &crpolr)

		client.On("GetRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&cr, nil)

		client.On("UpdateRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644, MatchType: "path", Type: "WAF", Name: "Test_Paths 3", Description: "AFW Test Extensions U", AverageThreshold: 5, BurstThreshold: 10, ClientIdentifier: "ip", UseXForwardForHeaders: true, RequestType: "ClientRequest", SameActionOnIpv6: false, Path: struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Values        []string "json:\"values\""
			}{PositiveMatch: true, Values: []string{"/login/", "/path/"}}, PathMatchType: "Custom", PathURIPositiveMatch: true, FileExtensions: struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Values        []string "json:\"values\""
			}{PositiveMatch: false, Values: []string{"3g2", "3gp", "aif", "aiff", "au", "avi", "bin", "bmp", "cab"}}, Hostnames: []string{"www.ludin.org"}, AdditionalMatchOptions: []struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Type          string   "json:\"type\""
				Values        []string "json:\"values\""
			}{struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Type          string   "json:\"type\""
				Values        []string "json:\"values\""
			}{PositiveMatch: true, Type: "IpAddressCondition", Values: []string{"198.129.76.39"}}, struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Type          string   "json:\"type\""
				Values        []string "json:\"values\""
			}{PositiveMatch: true, Type: "RequestMethodCondition", Values: []string{"GET"}}}, QueryParameters: []struct {
				Name          string   "json:\"name\""
				Values        []string "json:\"values\""
				PositiveMatch bool     "json:\"positiveMatch\""
				ValueInRange  bool     "json:\"valueInRange\""
			}{struct {
				Name          string   "json:\"name\""
				Values        []string "json:\"values\""
				PositiveMatch bool     "json:\"positiveMatch\""
				ValueInRange  bool     "json:\"valueInRange\""
			}{Name: "productId", Values: []string{"BUB_12", "SUSH_11"}, PositiveMatch: true, ValueInRange: false}}, CreateDate: "", UpdateDate: "", Used: false},
		).Return(&cu, nil)

		client.On("CreateRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, MatchType: "path", Type: "WAF", Name: "Test_Paths 3", Description: "AFW Test Extensions", AverageThreshold: 5, BurstThreshold: 10, ClientIdentifier: "ip", UseXForwardForHeaders: true, RequestType: "ClientRequest", SameActionOnIpv6: false, Path: struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Values        []string "json:\"values\""
			}{PositiveMatch: true, Values: []string{"/login/", "/path/"}}, PathMatchType: "Custom", PathURIPositiveMatch: true, FileExtensions: struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Values        []string "json:\"values\""
			}{PositiveMatch: false, Values: []string{"3g2", "3gp", "aif", "aiff", "au", "avi", "bin", "bmp", "cab"}}, Hostnames: []string{"www.ludin.org"}, AdditionalMatchOptions: []struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Type          string   "json:\"type\""
				Values        []string "json:\"values\""
			}{struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Type          string   "json:\"type\""
				Values        []string "json:\"values\""
			}{PositiveMatch: true, Type: "IpAddressCondition", Values: []string{"198.129.76.39"}}, struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				Type          string   "json:\"type\""
				Values        []string "json:\"values\""
			}{PositiveMatch: true, Type: "RequestMethodCondition", Values: []string{"GET"}}}, QueryParameters: []struct {
				Name          string   "json:\"name\""
				Values        []string "json:\"values\""
				PositiveMatch bool     "json:\"positiveMatch\""
				ValueInRange  bool     "json:\"valueInRange\""
			}{struct {
				Name          string   "json:\"name\""
				Values        []string "json:\"values\""
				PositiveMatch bool     "json:\"positiveMatch\""
				ValueInRange  bool     "json:\"valueInRange\""
			}{Name: "productId", Values: []string{"BUB_12", "SUSH_11"}, PositiveMatch: true, ValueInRange: false}}, CreateDate: "", UpdateDate: "", Used: false},
		).Return(&crpol, nil)

		client.On("RemoveRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&crpolr, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRatePolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "134644"),
						),
					},

					{
						Config: loadFixtureString("testdata/TestResRatePolicy/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "134644"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
