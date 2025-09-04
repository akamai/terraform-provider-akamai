package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiWAFProtection_res_basic(t *testing.T) {
	t.Run("match by WAFProtection ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		updateResponseAllProtectionsFalse := appsec.UpdateWAFProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFProtection/PolicyProtections.json"), &updateResponseAllProtectionsFalse)
		require.NoError(t, err)

		getResponseAllProtectionsFalse := appsec.GetWAFProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFProtection/PolicyProtections.json"), &getResponseAllProtectionsFalse)
		require.NoError(t, err)

		updateResponseOneProtectionTrue := appsec.UpdateWAFProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFProtection/UpdatedPolicyProtections.json"), &updateResponseOneProtectionTrue)
		require.NoError(t, err)

		getResponseOneProtectionTrue := appsec.GetWAFProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFProtection/UpdatedPolicyProtections.json"), &getResponseOneProtectionTrue)
		require.NoError(t, err)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.On("UpdateWAFProtection",
			testutils.MockContext,
			appsec.UpdateWAFProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()
		client.On("GetWAFProtection",
			testutils.MockContext,
			appsec.GetWAFProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.On("GetWAFProtection",
			testutils.MockContext,
			appsec.GetWAFProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Delete, performed automatically to clean up
		client.On("UpdateWAFProtection",
			testutils.MockContext,
			appsec.UpdateWAFProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResWAFProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_waf_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_waf_protection.test", "enabled", "false"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
