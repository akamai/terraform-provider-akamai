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

		client.On("GetConfigurationClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationCloneRequest{ConfigID: 43253, Version: 15},
		).Return(&cr, nil)

		client.On("CreateConfigurationClone",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateConfigurationCloneRequest{ConfigID: 43253, CreateFromVersion: 7},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfigurationClone/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_version_clone.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
