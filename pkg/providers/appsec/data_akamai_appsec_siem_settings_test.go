package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSiemSettings_data_basic(t *testing.T) {
	t.Run("match by SiemSettings ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getSiemSettingsResponse := appsec.GetSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSSiemSettings/SiemSettings.json"), &getSiemSettingsResponse)
		require.NoError(t, err)

		client.On("GetSiemSettings",
			testutils.MockContext,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSSiemSettings/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_siem_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
