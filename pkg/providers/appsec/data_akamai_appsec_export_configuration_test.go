package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiExportConfiguration_data_basic(t *testing.T) {
	t.Run("match by ExportConfiguration ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetExportConfigurationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSExportConfiguration/ExportConfiguration.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetExportConfigurations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetExportConfigurationsRequest{ConfigID: 43253, Version: 7},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSExportConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_export_configuration.test", "id", "0"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
