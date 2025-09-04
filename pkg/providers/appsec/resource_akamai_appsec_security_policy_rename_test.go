package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSecurityPolicyRename_res_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateSecurityPolicyResponse := appsec.UpdateSecurityPolicyResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyRename/SecurityPolicyUpdate.json"), &updateSecurityPolicyResponse)
		require.NoError(t, err)

		getSecurityPolicyResponse := appsec.GetSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyRename/SecurityPolicy.json"), &getSecurityPolicyResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSecurityPolicy",
			testutils.MockContext,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049"},
		).Return(&getSecurityPolicyResponse, nil)

		client.On("UpdateSecurityPolicy",
			testutils.MockContext,
			appsec.UpdateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLE_114049", PolicyName: "Cloned Test for Launchpad 15"},
		).Return(&updateSecurityPolicyResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyRename/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_rename.test", "id", "43253:PLE_114049"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyRename/update_by_id.tf"),
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
