package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSecurityPolicyDefaultProtections_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

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

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetSecurityPolicy",
			testutils.MockContext,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049"},
		).Return(&getSecurityPolicyResponse, nil).Times(3)

		client.APPSEC.On("GetSecurityPolicy",
			testutils.MockContext,
			appsec.GetSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049"},
		).Return(&getSecurityPolicyAfterUpdateResponse, nil).Twice()

		client.APPSEC.On("UpdateSecurityPolicy",
			testutils.MockContext,
			appsec.UpdateSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049", PolicyName: "PLEB Cloned Test for Launchpad 15 New"},
		).Return(&updateSecurityPolicyResponse, nil)

		client.APPSEC.On("CreateSecurityPolicyWithDefaultProtections",
			testutils.MockContext,
			appsec.CreateSecurityPolicyWithDefaultProtectionsRequest{ConfigVersion: appsec.ConfigVersion{ConfigID: 43253, Version: 7}, PolicyName: "PLEB Cloned Test for Launchpad 15", PolicyPrefix: "PLEB"},
		).Return(&createSecurityPolicyResponse, nil)

		client.APPSEC.On("RemoveSecurityPolicy",
			testutils.MockContext,
			appsec.RemoveSecurityPolicyRequest{ConfigID: 43253, Version: 7, PolicyID: "PLEB_114049"},
		).Return(&removeSecurityPolicyResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
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

		client.APPSEC.AssertExpectations(t)
	})

}

func TestAkamaiSecurityPolicyDefaultProtections_res_failure_creating_policy(t *testing.T) {
	t.Parallel()
	t.Run("match by SecurityPolicy ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("CreateSecurityPolicyWithDefaultProtections",
			testutils.MockContext,
			appsec.CreateSecurityPolicyWithDefaultProtectionsRequest{ConfigVersion: appsec.ConfigVersion{ConfigID: 43253, Version: 7}, PolicyName: "PLEB Cloned Test for Launchpad 15", PolicyPrefix: "PLEB"},
		).Return(nil, fmt.Errorf("create security policy request failed: policy name already in use"))

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
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

		client.APPSEC.AssertExpectations(t)
	})
}
