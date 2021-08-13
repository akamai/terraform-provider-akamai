package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSelectedHostname_res_basic(t *testing.T) {
	t.Run("match by SelectedHostname ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSelectedHostnameResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSelectedHostnameResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		hns := appsec.GetSelectedHostnameResponse{}
		expectJSHN := compactJSON(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"))
		json.Unmarshal([]byte(expectJSHN), &hns)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSelectedHostname",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSelectedHostnameRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateSelectedHostname",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSelectedHostnameRequest{ConfigID: 43253, Version: 7, HostnameList: []appsec.Hostname{
				{
					Hostname: "rinaldi.sandbox.akamaideveloper.com",
				},
				{
					Hostname: "sujala.sandbox.akamaideveloper.com",
				},
			},
			},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSelectedHostname/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_selected_hostnames.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
