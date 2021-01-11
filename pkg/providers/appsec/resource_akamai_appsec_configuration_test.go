package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfiguration_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateConfigurationResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/ConfigurationUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetConfigurationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveConfigurationResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		ccr := appsec.CreateConfigurationResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/ConfigurationCreate.json"))
		json.Unmarshal([]byte(expectJSC), &ccr)

		client.On("CreateConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateConfigurationRequest{Name: "Akamai Tools New", Description: "TF Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&ccr, nil)

		client.On("GetConfigurations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationsRequest{ConfigID: 432531, Name: "Akamai Tools New"},
		).Return(&cr, nil)

		client.On("UpdateConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateConfigurationRequest{ConfigID: 432531, Name: "Akamai Tools New", Description: "TF Tools 1"},
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
						Config: loadFixtureString("testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "432531"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResConfiguration/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "432531"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
