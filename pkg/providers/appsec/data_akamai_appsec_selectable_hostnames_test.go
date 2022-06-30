package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSelectableHostnames_data_basic(t *testing.T) {
	t.Run("match by SelectableHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getSelectableHostnamesResponse := appsec.GetSelectableHostnamesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSSelectableHostnames/SelectableHostnames.json"), &getSelectableHostnamesResponse)

		client.On("GetSelectableHostnames",
			mock.Anything,
			appsec.GetSelectableHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectableHostnamesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSSelectableHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_selectable_hostnames.test", "id", "0"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
