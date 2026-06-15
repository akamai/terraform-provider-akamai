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

func TestAkamaiAdvancedSettingsRequestBodyDataBasic(t *testing.T) {
	t.Parallel()
	t.Run("match by AdvancedSettingsRequestBody ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRequestBodyResponse := appsec.GetAdvancedSettingsRequestBodyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"), &getRequestBodyResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetAdvancedSettingsRequestBody",
			testutils.MockContext,
			appsec.GetAdvancedSettingsRequestBodyRequest{ConfigID: 43253, Version: 7},
		).Return(&getRequestBodyResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsRequestBody/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_request_body.test", "id", "43253:"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}

func TestAkamaiAdvancedSettingsRequestBodyDataBasicPolicyID(t *testing.T) {
	t.Parallel()
	t.Run("match by AdvancedSettingsRequestBodyPolicy ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRequestBodyResponse := appsec.GetAdvancedSettingsRequestBodyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json"), &getRequestBodyResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetAdvancedSettingsRequestBody",
			testutils.MockContext,
			appsec.GetAdvancedSettingsRequestBodyRequest{ConfigID: 43253, Version: 7, PolicyID: "test_policy"},
		).Return(&getRequestBodyResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsRequestBody/match_by_policy_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_request_body.policy", "id", "43253:test_policy"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
