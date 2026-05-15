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

func TestAkamaiAPICoProtection_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by APIConstraintsProtection ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		updateResponseAllProtectionsFalse := appsec.UpdateAPIConstraintsProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAPIConstraintsProtection/PolicyProtections.json"), &updateResponseAllProtectionsFalse)
		require.NoError(t, err)

		getResponseAllProtectionsFalse := appsec.GetAPIConstraintsProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAPIConstraintsProtection/PolicyProtections.json"), &getResponseAllProtectionsFalse)
		require.NoError(t, err)

		updateResponseOneProtectionTrue := appsec.UpdateAPIConstraintsProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAPIConstraintsProtection/UpdatedPolicyProtections.json"), &updateResponseOneProtectionTrue)
		require.NoError(t, err)

		getResponseOneProtectionTrue := appsec.GetAPIConstraintsProtectionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAPIConstraintsProtection/UpdatedPolicyProtections.json"), &getResponseOneProtectionTrue)
		require.NoError(t, err)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.APPSEC.On("UpdateAPIConstraintsProtection",
			testutils.MockContext,
			appsec.UpdateAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()
		client.APPSEC.On("GetAPIConstraintsProtection",
			testutils.MockContext,
			appsec.GetAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.APPSEC.On("GetAPIConstraintsProtection",
			testutils.MockContext,
			appsec.GetAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()
		client.APPSEC.On("GetAPIConstraintsProtection",
			testutils.MockContext,
			appsec.GetAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Update, with terminal Read
		client.APPSEC.On("UpdateAPIConstraintsProtection",
			testutils.MockContext,
			appsec.UpdateAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230",
				ApplyAPIConstraints: true},
		).Return(&updateResponseOneProtectionTrue, nil).Once()
		client.APPSEC.On("GetAPIConstraintsProtection",
			testutils.MockContext,
			appsec.GetAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Read, performed as part of "id" check.
		client.APPSEC.On("GetAPIConstraintsProtection",
			testutils.MockContext,
			appsec.GetAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Delete, performed automatically to clean up
		client.APPSEC.On("UpdateAPIConstraintsProtection",
			testutils.MockContext,
			appsec.UpdateAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAPIConstraintsProtection/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "enabled", "false"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAPIConstraintsProtection/update_by_id.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "enabled", "true"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
