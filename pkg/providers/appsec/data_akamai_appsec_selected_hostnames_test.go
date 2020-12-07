package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSelectedHostnames_data_basic(t *testing.T) {
	t.Run("match by SelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetSelectedHostnamesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSSelectedHostnames/SelectedHostnames.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetSelectedHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_selected_hostnames.test", "id", "43253:7"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
