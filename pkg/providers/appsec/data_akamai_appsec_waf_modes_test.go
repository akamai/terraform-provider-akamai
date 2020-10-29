package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiWAFModes_data_basic(t *testing.T) {
	t.Run("match by WAFModes ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetWAFModesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSWAFModes/WAFModes.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetWAFModes",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetWAFModesRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAFModes/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_waf_modes.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
