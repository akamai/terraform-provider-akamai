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

		cr := appsec.GetSecurityPolicyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crp := appsec.CreateSecurityPolicyResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicyCreate.json"))
		json.Unmarshal([]byte(expectJSC), &crp)

		rp := appsec.RemoveSecurityPolicyResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"))
		json.Unmarshal([]byte(expectJSR), &rp)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&cr, nil)

		client.On("CreateSecurityPolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyName: "Cloned Test for Launchpad 15", PolicyPrefix: "LN", DefaultSettings: true},
		).Return(&crp, nil)

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
							resource.TestCheckResourceAttr("akamai_appsec_security_policy.test", "id", "43253:PLE_114049"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
