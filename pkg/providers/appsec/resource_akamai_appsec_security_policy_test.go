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

		getSecurityPolicyResponse := appsec.GetSecurityPolicyResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"), &getSecurityPolicyResponse)

		createSecurityPolicyResponse := appsec.CreateSecurityPolicyResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicyCreate.json"), &createSecurityPolicyResponse)

		removeSecurityPolicyResponse := appsec.RemoveSecurityPolicyResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResSecurityPolicy/SecurityPolicy.json"), &removeSecurityPolicyResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicy",
			mock.Anything,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&getSecurityPolicyResponse, nil)

		client.On("CreateSecurityPolicy",
			mock.Anything,
			appsec.CreateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyName: "PLE Cloned Test for Launchpad 15", PolicyPrefix: "PLE", DefaultSettings: true},
		).Return(&createSecurityPolicyResponse, nil)

		client.On("RemoveSecurityPolicy",
			mock.Anything,
			appsec.RemoveSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&removeSecurityPolicyResponse, nil)

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
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
