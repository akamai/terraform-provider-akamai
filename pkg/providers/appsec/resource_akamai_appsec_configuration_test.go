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
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetConfigurationResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationRequest{},
		).Return(&cr, nil)

		client.On("UpdateConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateConfigurationRequest{},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
