package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomDeny_res_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateCustomDenyResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetCustomDenyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveCustomDenyResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("RemoveCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&crd, nil)

		crc := appsec.CreateCustomDenyResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyCreate.json"))
		json.Unmarshal([]byte(expectJSC), &crc)

		client.On("GetCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&cr, nil)

		client.On("UpdateCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, Description: "testing", Name: "new_custom_deny", ID: "deny_custom_622918", Parameters: []struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "response_status_code", Value: "403"}, struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "prevent_browser_cache", Value: "false"}, struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "response_content_type", Value: "application/json"}, struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "response_body_content", Value: "new testing"}}},
		).Return(&cu, nil)

		client.On("CreateCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, Description: "testing", Name: "new_custom_deny", Parameters: []struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "response_status_code", Value: "403"}, struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "prevent_browser_cache", Value: "true"}, struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "response_content_type", Value: "application/json"}, struct {
				DisplayName string "json:\"displayName\""
				Name        string "json:\"name\""
				Value       string "json:\"value\""
			}{DisplayName: "", Name: "response_body_content", Value: "new testing"}}},
		).Return(&crc, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "deny_custom_622918"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "deny_custom_622918"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
