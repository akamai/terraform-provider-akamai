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
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetCustomDenyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&cr, nil)

		client.On("UpdateCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
