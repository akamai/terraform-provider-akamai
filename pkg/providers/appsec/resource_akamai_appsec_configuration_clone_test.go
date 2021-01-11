package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfigurationClone_res_basic(t *testing.T) {
	t.Run("match by ConfigurationClone ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.CreateConfigurationCloneResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResConfigurationClone/ConfigurationCloneCreated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetConfigurationCloneResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResConfigurationClone/ConfigurationClone.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveConfigurationResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResConfigurationClone/ConfigurationClone.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetConfigurationClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationCloneRequest{ConfigID: 43253, Version: 15},
		).Return(&cr, nil)

		client.On("CreateConfigurationClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateConfigurationCloneRequest{Name: "Test Configuratin", Description: "New configuration test", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}, CreateFrom: struct {
				ConfigID int "json:\"configId\""
				Version  int "json:\"version\""
			}{ConfigID: 43253, Version: 7}},
		).Return(&cu, nil)

		client.On("RemoveConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfigurationClone/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_clone.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
