package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSiemSettings_res_basic(t *testing.T) {
	t.Run("match by SiemSettings ID", func(t *testing.T) {
		client := &appsec.Mock{}

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

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil)

		client.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: ptr.To(true), SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"}, Exceptions: []appsec.Exception{}},
		).Return(&updateSiemSettingsResponse, nil)

		client.On("RemoveSiemSettings",
			testutils.MockContext,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, FirewallPolicyIDs: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by SiemSettings ID when SIEM exceptions are added", func(t *testing.T) {
		client := &appsec.Mock{}

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

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil)

		client.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"},
				Exceptions: []appsec.Exception{
					{
						Protection:  "rate",
						ActionTypes: []string{"alert"},
					},
				}},
		).Return(&updateSiemSettingsResponse, nil)

		client.On("RemoveSiemSettings",
			testutils.MockContext,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: false, FirewallPolicyIDs: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id_exceptions_enabled.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by SiemSettings ID when SIEM exceptions are added with empty actions", func(t *testing.T) {
		client := &appsec.Mock{}

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("update after removing exceptions block", func(t *testing.T) {
		client := &appsec.Mock{}

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

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil).Times(8)

		client.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil).Times(2)

		client.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"},
				Exceptions: []appsec.Exception{
					{
						Protection:  "rate",
						ActionTypes: []string{"alert"},
					},
				}},
		).Return(&updateSiemSettingsResponse, nil).Times(1)

		client.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil).Times(1)

		client.On("UpdateSiemSettings",
			testutils.MockContext,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: ptr.To(true), SiemDefinitionID: 1, FirewallPolicyIDs: []string{"12345"}, Exceptions: []appsec.Exception{}},
		).Return(&updateSiemSettingsResponseNoExceptions, nil).Times(1)

		client.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponseNoExceptions, nil).Times(2)

		client.On("RemoveSiemSettings",
			testutils.MockContext,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: false, FirewallPolicyIDs: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("match by SiemSettings ID when exceptions block is empty", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id_empty_exceptions_block.tf"),
						ExpectError: regexp.MustCompile(`Error: Invalid exceptions configuration`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
