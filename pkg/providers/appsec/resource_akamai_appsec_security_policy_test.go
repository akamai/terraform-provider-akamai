package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSecurityPolicy_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getSecurityPolicyResponse := appsec.GetSecurityPolicyResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicy/SecurityPolicy.json"), &getSecurityPolicyResponse)
		require.NoError(t, err)

		createSecurityPolicyResponse := appsec.CreateSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicy/SecurityPolicyCreate.json"), &createSecurityPolicyResponse)
		require.NoError(t, err)

		removeSecurityPolicyResponse := appsec.RemoveSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicy/SecurityPolicy.json"), &removeSecurityPolicyResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetSecurityPolicy",
			testutils.MockContext,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&getSecurityPolicyResponse, nil)

		client.APPSEC.On("CreateSecurityPolicy",
			testutils.MockContext,
			appsec.CreateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyName: "PLE Cloned Test for Launchpad 15", PolicyPrefix: "PLE", DefaultSettings: true},
		).Return(&createSecurityPolicyResponse, nil)

		client.APPSEC.On("RemoveSecurityPolicy",
			testutils.MockContext,
			appsec.RemoveSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&removeSecurityPolicyResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicy/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_security_policy.test", "id", "43253:PLE_114049"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
