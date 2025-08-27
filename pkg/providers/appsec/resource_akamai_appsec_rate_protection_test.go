package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRateProtection_res_basic(t *testing.T) {
	t.Run("match by RateProtection ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		updateResponseAllProtectionsFalse := appsec.UpdateRateProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRateProtection/PolicyProtections.json"), &updateResponseAllProtectionsFalse)
		require.NoError(t, err)

		getResponseAllProtectionsFalse := appsec.GetRateProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRateProtection/PolicyProtections.json"), &getResponseAllProtectionsFalse)
		require.NoError(t, err)

		updateResponseOneProtectionTrue := appsec.UpdateRateProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRateProtection/UpdatedPolicyProtections.json"), &updateResponseOneProtectionTrue)
		require.NoError(t, err)

		getResponseOneProtectionTrue := appsec.GetRateProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRateProtection/UpdatedPolicyProtections.json"), &getResponseOneProtectionTrue)
		require.NoError(t, err)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.On("UpdateRateProtection",
			testutils.MockContext,
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()
		client.On("GetRateProtection",
			testutils.MockContext,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.On("GetRateProtection",
			testutils.MockContext,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()
		client.On("GetRateProtection",
			testutils.MockContext,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Update, with terminal Read
		client.On("UpdateRateProtection",
			testutils.MockContext,
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230",
				ApplyRateControls: true},
		).Return(&updateResponseOneProtectionTrue, nil).Once()
		client.On("GetRateProtection",
			testutils.MockContext,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Read, performed as part of "id" check.
		client.On("GetRateProtection",
			testutils.MockContext,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Delete, performed automatically to clean up
		client.On("UpdateRateProtection",
			testutils.MockContext,
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRateProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "enabled", "false"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRateProtection/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "enabled", "true"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
