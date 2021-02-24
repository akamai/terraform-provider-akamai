package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfigurationVersionClone_res_basic(t *testing.T) {
	t.Run("match by ConfigurationVersionClone ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.CreateConfigurationVersionCloneResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResConfigurationVersionClone/ConfigurationVersionClone.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetConfigurationVersionCloneResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResConfigurationVersionClone/ConfigurationVersionClone.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveConfigurationVersionCloneResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResConfigurationVersionClone/ConfigurationVersionClone.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetConfigurationVersionClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationVersionCloneRequest{ConfigID: 43253, Version: 15},
		).Return(&cr, nil)

		client.On("CreateConfigurationVersionClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateConfigurationVersionCloneRequest{ConfigID: 43253, CreateFromVersion: 7},
		).Return(&cu, nil)

		client.On("RemoveConfigurationVersionClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveConfigurationVersionCloneRequest{ConfigID: 43253, Version: 15},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfigurationVersionClone/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_version_clone.test", "id", "15"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
