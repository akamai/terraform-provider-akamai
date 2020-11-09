package appsec

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSecurityPolicyClone_res_basic(t *testing.T) {
	t.Run("match by SecurityPolicyClone ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.CreateSecurityPolicyCloneResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicyClone/SecurityPolicyCloneCreated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSecurityPolicyCloneResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicyClone/SecurityPolicyClone.json"))
		json.Unmarshal([]byte(expectJS), &cr)
		fmt.Sprintf("TEST %v", cr)
		client.On("GetSecurityPolicyClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSecurityPolicyCloneRequest{ConfigID: 43253, Version: 15, PolicyID: "LNPD_76189"},
		).Return(&cr, nil)

		client.On("CreateSecurityPolicyClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateSecurityPolicyCloneRequest{ConfigID: 43253, Version: 15, CreateFromSecurityPolicy: "LNPD_76189", PolicyName: "Cloned Test for Launchpad 15", PolicyPrefix: "LN"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicyClone/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_clone.test", "id", "LNPD_76189"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
