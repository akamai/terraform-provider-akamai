package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiMatchTargets_data_basic(t *testing.T) {
	t.Run("match by MatchTargets ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetMatchTargetsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSMatchTargets/MatchTargets.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetMatchTargets",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetMatchTargetsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSMatchTargets/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_match_targets.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
