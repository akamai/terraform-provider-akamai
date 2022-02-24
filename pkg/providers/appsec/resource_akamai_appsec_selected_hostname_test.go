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

		cu := appsec.UpdateSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json")), &cu)

		cr := appsec.GetSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json")), &cr)

		hns := appsec.GetSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json")), &hns)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSelectedHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateSelectedHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSelectedHostnamesRequest{ConfigID: 43253, Version: 7, HostnameList: []appsec.Hostname{
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
