package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSiemSettings_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by SiemSettings ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateSiemSettingsResponse := appsec.UpdateSiemSettingsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &updateSiemSettingsResponse)
		require.NoError(t, err)

		getSiemSettingsResponse := appsec.GetSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &getSiemSettingsResponse)
		require.NoError(t, err)

		removeSiemSettingsResponse := appsec.RemoveSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &removeSiemSettingsResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil)

		client.APPSEC.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: ptr.To(true), IncludeJA4FingerprintToSiem: ptr.To(true), SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"}, Exceptions: []appsec.Exception{}},
		).Return(&updateSiemSettingsResponse, nil)

		client.APPSEC.On("RemoveSiemSettings",
			testutils.MockContext,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, FirewallPolicyIDs: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by SiemSettings ID when SIEM exceptions are added", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateSiemSettingsResponse := appsec.UpdateSiemSettingsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettingsExceptionsEnabled.json"), &updateSiemSettingsResponse)
		require.NoError(t, err)

		getSiemSettingsResponse := appsec.GetSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettingsExceptionsEnabled.json"), &getSiemSettingsResponse)
		require.NoError(t, err)

		removeSiemSettingsResponse := appsec.RemoveSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettingsExceptionsEnabled.json"), &removeSiemSettingsResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil)

		client.APPSEC.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"},
				Exceptions: []appsec.Exception{
					{
						Protection:  "rate",
						ActionTypes: []string{"alert"},
					},
				}},
		).Return(&updateSiemSettingsResponse, nil)

		client.APPSEC.On("RemoveSiemSettings",
			testutils.MockContext,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: false, FirewallPolicyIDs: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id_exceptions_enabled.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by SiemSettings ID when SIEM exceptions are added with empty actions", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id_exceptions_enabled_empty_input.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
					),
					ExpectError: regexp.MustCompile(`Error: Not enough list items`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("update after removing exceptions block", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateSiemSettingsResponse := appsec.UpdateSiemSettingsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettingsExceptionsEnabled.json"), &updateSiemSettingsResponse)
		require.NoError(t, err)

		updateSiemSettingsResponseNoExceptions := appsec.UpdateSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &updateSiemSettingsResponseNoExceptions)
		require.NoError(t, err)

		getSiemSettingsResponse := appsec.GetSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettingsExceptionsEnabled.json"), &getSiemSettingsResponse)
		require.NoError(t, err)

		getSiemSettingsResponseNoExceptions := appsec.GetSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &getSiemSettingsResponseNoExceptions)
		require.NoError(t, err)

		removeSiemSettingsResponse := appsec.RemoveSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettingsExceptionsEnabled.json"), &removeSiemSettingsResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil).Times(8)

		client.APPSEC.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil).Times(2)

		client.APPSEC.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"},
				Exceptions: []appsec.Exception{
					{
						Protection:  "rate",
						ActionTypes: []string{"alert"},
					},
				}},
		).Return(&updateSiemSettingsResponse, nil).Times(1)

		client.APPSEC.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil).Times(1)

		client.APPSEC.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: ptr.To(true), IncludeJA4FingerprintToSiem: ptr.To(true), SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"}, Exceptions: []appsec.Exception{}},
		).Return(&updateSiemSettingsResponseNoExceptions, nil).Times(1)

		client.APPSEC.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponseNoExceptions, nil).Times(2)

		client.APPSEC.On("RemoveSiemSettings",
			testutils.MockContext,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: false, FirewallPolicyIDs: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id_exceptions_enabled.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by SiemSettings ID when exceptions block is empty", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id_empty_exceptions_block.tf"),
					ExpectError: regexp.MustCompile(`Error: invalid exceptions configuration`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
