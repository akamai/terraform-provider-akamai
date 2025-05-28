package appsec

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSecurityPolicy_data_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		securityPoliciesBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSSecurityPolicy/SecurityPolicy.json")
		getSecurityPoliciesResponse := appsec.GetSecurityPoliciesResponse{}
		err := json.Unmarshal(securityPoliciesBytes, &getSecurityPoliciesResponse)
		require.NoError(t, err)

		securityPoliciesJSONBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSSecurityPolicy/SecurityPolicyJSON.json")
		buf := &bytes.Buffer{}
		err = json.Compact(buf, securityPoliciesJSONBytes)
		require.NoError(t, err)
		securityPoliciesJSONString := buf.String()

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicies",
			testutils.MockContext,
			appsec.GetSecurityPoliciesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSecurityPoliciesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSSecurityPolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_security_policy.test", "id", "43253:7"),
							resource.TestCheckResourceAttr("data.akamai_appsec_security_policy.test", "json", securityPoliciesJSONString),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
