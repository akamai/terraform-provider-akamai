package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiSecurityPolicy_data_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &mockappsec{}

		getSecurityPoliciesResponse := appsec.GetSecurityPoliciesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSSecurityPolicy/SecurityPolicy.json"), &getSecurityPoliciesResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicies",
			mock.Anything,
			appsec.GetSecurityPoliciesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSecurityPoliciesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSSecurityPolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_security_policy.test", "id", "43253:7"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
