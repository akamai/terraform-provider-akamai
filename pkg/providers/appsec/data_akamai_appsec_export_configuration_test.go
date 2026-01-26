package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiExportConfiguration_data_basic(t *testing.T) {
	t.Run("Configuration Export Tests", func(t *testing.T) {
		client := &appsec.Mock{}

		getExportConfigurationResponse := appsec.GetExportConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSExportConfiguration/ExportConfiguration.json"), &getExportConfigurationResponse)
		require.NoError(t, err)

		client.On("GetExportConfiguration",
			testutils.MockContext,
			appsec.GetExportConfigurationRequest{ConfigID: 43253, Version: 7},
		).Return(&getExportConfigurationResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSExportConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_export_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
