package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSelectedHostname_res_basic(t *testing.T) {
	t.Run("match by SelectedHostname ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateSelectedHostnamesResponse := appsec.UpdateSelectedHostnamesResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"), &updateSelectedHostnamesResponse)
		require.NoError(t, err)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		getSelectedHostnamesResponseAfterUpdate := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"), &getSelectedHostnamesResponseAfterUpdate)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSelectedHostnames",
			mock.Anything,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.On("UpdateSelectedHostnames",
			mock.Anything,
			appsec.UpdateSelectedHostnamesRequest{ConfigID: 43253, Version: 7, HostnameList: []appsec.Hostname{
				{
					Hostname: "rinaldi.sandbox.akamaideveloper.com",
				},
				{
					Hostname: "sujala.sandbox.akamaideveloper.com",
				},
			},
			},
		).Return(&updateSelectedHostnamesResponse, nil)

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
