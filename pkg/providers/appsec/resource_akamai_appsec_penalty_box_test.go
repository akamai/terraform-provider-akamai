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
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResPenaltyBox/PenaltyBox.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetPenaltyBoxResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResPenaltyBox/PenaltyBox.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetPenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdatePenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", PenaltyBoxProtection: true},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResPenaltyBox/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_penalty_box.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
