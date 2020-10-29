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

		cv := appsec.GetSelectableHostnamesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSSelectableHostnames/SelectableHostnames.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetSelectableHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSelectableHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&cv, nil)

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
