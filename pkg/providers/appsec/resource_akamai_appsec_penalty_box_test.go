package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiPenaltyBox_res_basic(t *testing.T) {
	t.Run("match by PenaltyBox ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdatePenaltyBoxResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResPenaltyBox/PenaltyBox.json")), &cu)

		cr := appsec.GetPenaltyBoxResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResPenaltyBox/PenaltyBox.json")), &cr)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetPenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdatePenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", PenaltyBoxProtection: false},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResPenaltyBox/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_penalty_box.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
