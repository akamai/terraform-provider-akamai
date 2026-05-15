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

func TestAkamaiAdvancedSettingsPragma_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by AdvancedSettingsPragma ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getPragmaResponse := appsec.GetAdvancedSettingsPragmaResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &getPragmaResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetAdvancedSettingsPragma",
			testutils.MockContext,
			appsec.GetAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7},
		).Return(&getPragmaResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsPragma/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_pragma_header.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
