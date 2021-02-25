package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSecurityPolicyRename_res_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSecurityPolicyResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicyRename/SecurityPolicyUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSecurityPolicyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicyRename/SecurityPolicy.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&cr, nil)

		client.On("UpdateSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049", PolicyName: "Cloned Test for Launchpad 15", DefaultSettings: false, PolicyPrefix: ""},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicyRename/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_rename.test", "id", "43253:7:PLE_114049"),
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicyRename/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_rename.test", "id", "43253:7:PLE_114049"),
						),
						//ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
