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

		updatePenaltyBoxResponse := appsec.UpdatePenaltyBoxResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResPenaltyBox/PenaltyBox.json"), &updatePenaltyBoxResponse)

		getPenaltyBoxResponse := appsec.GetPenaltyBoxResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResPenaltyBox/PenaltyBox.json"), &getPenaltyBoxResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetPenaltyBox",
			mock.Anything,
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getPenaltyBoxResponse, nil)

		client.On("UpdatePenaltyBox",
			mock.Anything,
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", PenaltyBoxProtection: false},
		).Return(&updatePenaltyBoxResponse, nil)

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
