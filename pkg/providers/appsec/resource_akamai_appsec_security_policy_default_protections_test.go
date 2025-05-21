package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSecurityPolicyDefaultProtections_res_basic(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getSecurityPolicyResponse := appsec.GetSecurityPolicyResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyDefaultProtections/SecurityPolicy.json"), &getSecurityPolicyResponse)
		require.NoError(t, err)

		getSecurityPolicyAfterUpdateResponse := appsec.GetSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyDefaultProtections/SecurityPolicyDefaultProtectionsUpdated.json"), &getSecurityPolicyAfterUpdateResponse)
		require.NoError(t, err)

		createSecurityPolicyResponse := appsec.CreateSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyDefaultProtections/SecurityPolicyDefaultProtectionsCreate.json"), &createSecurityPolicyResponse)
		require.NoError(t, err)

		updateSecurityPolicyResponse := appsec.UpdateSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyDefaultProtections/SecurityPolicyDefaultProtectionsUpdated.json"), &updateSecurityPolicyResponse)
		require.NoError(t, err)

		removeSecurityPolicyResponse := appsec.RemoveSecurityPolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSecurityPolicyDefaultProtections/SecurityPolicy.json"), &removeSecurityPolicyResponse)
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
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049"},
		).Return(&getSecurityPolicyResponse, nil).Times(3)

		client.On("GetSecurityPolicy",
			testutils.MockContext,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049"},
		).Return(&getSecurityPolicyAfterUpdateResponse, nil).Twice()

		client.On("UpdateSecurityPolicy",
			testutils.MockContext,
			appsec.UpdateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049", PolicyName: "PLEB Cloned Test for Launchpad 15 New"},
		).Return(&updateSecurityPolicyResponse, nil)

		client.On("CreateSecurityPolicyWithDefaultProtections",
			testutils.MockContext,
			appsec.CreateSecurityPolicyWithDefaultProtectionsRequest{ConfigVersion: appsec.ConfigVersion{ConfigID: 43253, Version: 7}, PolicyName: "PLEB Cloned Test for Launchpad 15", PolicyPrefix: "PLEB"},
		).Return(&createSecurityPolicyResponse, nil)

		client.On("RemoveSecurityPolicy",
			testutils.MockContext,
			appsec.RemoveSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049"},
		).Return(&removeSecurityPolicyResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyDefaultProtections/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_default_protections.test", "id", "43253:PLEB_114049"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyDefaultProtections/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_default_protections.test", "security_policy_name", "PLEB Cloned Test for Launchpad 15 New"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiSecurityPolicyDefaultProtections_res_failure_creating_policy(t *testing.T) {
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("CreateSecurityPolicyWithDefaultProtections",
			testutils.MockContext,
			appsec.CreateSecurityPolicyWithDefaultProtectionsRequest{ConfigVersion: appsec.ConfigVersion{ConfigID: 43253, Version: 7}, PolicyName: "PLEB Cloned Test for Launchpad 15", PolicyPrefix: "PLEB"},
		).Return(nil, fmt.Errorf("create security policy request failed: policy name already in use"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyDefaultProtections/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_default_protections.test", "id", "43253:PLEB_114049"),
						),
						ExpectError: regexp.MustCompile(`create security policy request failed: policy name already in use`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
