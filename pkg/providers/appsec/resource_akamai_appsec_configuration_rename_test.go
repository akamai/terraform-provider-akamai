package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfigurationRename_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateConfigurationResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResConfigurationRename/ConfigurationUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetConfigurationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResConfigurationRename/Configuration.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveConfigurationResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResConfigurationRename/Configuration.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetConfigurations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationsRequest{ConfigID: 432531, Name: "Akamai Tools New"},
		).Return(&cr, nil)

		client.On("UpdateConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateConfigurationRequest{ConfigID: 432531, Name: "Akamai Tools New", Description: "TF Tools"},
		).Return(&cu, nil)

		client.On("RemoveConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveConfigurationRequest{ConfigID: 432531},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfigurationRename/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_rename.test", "id", "432531"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResConfigurationRename/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_rename.test", "id", "432531"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
