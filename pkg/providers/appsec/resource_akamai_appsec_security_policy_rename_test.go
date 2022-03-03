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
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSecurityPolicyRename/SecurityPolicyUpdate.json")), &cu)

		cr := appsec.GetSecurityPolicyResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSecurityPolicyRename/SecurityPolicy.json")), &cr)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&cr, nil)

		client.On("UpdateSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049", PolicyName: "Cloned Test for Launchpad 15"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicyRename/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_rename.test", "id", "43253:PLE_114049"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResSecurityPolicyRename/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_rename.test", "id", "43253:PLE_114049"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
