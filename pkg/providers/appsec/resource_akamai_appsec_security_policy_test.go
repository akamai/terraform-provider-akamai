package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSecurityPolicy_res_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSecurityPolicyResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicyUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSecurityPolicyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crp := appsec.CreateSecurityPolicyResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicyCreate.json"))
		json.Unmarshal([]byte(expectJSC), &crp)

		rp := appsec.RemoveSecurityPolicyResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"))
		json.Unmarshal([]byte(expectJSR), &rp)

		client.On("GetSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&cr, nil)

		client.On("CreateSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyName: "Cloned Test for Launchpad 15", PolicyPrefix: "LN"},
		).Return(&crp, nil)

		client.On("UpdateSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049", PolicyName: "Cloned Test for Launchpad 21", PolicyPrefix: "LN"},
		).Return(&cu, nil)

		client.On("RemoveSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&rp, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy.test", "id", "PLE_114049"),
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicy/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy.test", "id", "PLE_114049"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
