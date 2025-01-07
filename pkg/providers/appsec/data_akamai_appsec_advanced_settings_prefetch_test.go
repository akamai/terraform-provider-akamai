package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsPrefetch_data_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPrefetch ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getPrefetchResponse := appsec.GetAdvancedSettingsPrefetchResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"), &getPrefetchResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsPrefetch",
			testutils.MockContext,
			appsec.GetAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7},
		).Return(&getPrefetchResponse, nil)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
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
						Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsPrefetch/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_prefetch.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
