package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiMatchTargetSequence_res_basic(t *testing.T) {
	t.Run("match by MatchTargetSequence ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateMatchTargetSequenceResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResMatchTargetSequence/MatchTargetSequenceResp.json"))

		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetMatchTargetSequenceResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResMatchTargetSequence/MatchTargetSequence.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		mts := appsec.GetMatchTargetSequencesResponse{}
		expectJSMTS := compactJSON(loadFixtureBytes("testdata/TestResMatchTargetSequence/MatchTargetSequences.json"))
		json.Unmarshal([]byte(expectJSMTS), &mts)

		client.On("GetMatchTargetSequence",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetMatchTargetSequenceRequest{ConfigID: 43253, ConfigVersion: 7, Type: "website"},
		).Return(&cr, nil)

		client.On("UpdateMatchTargetSequence",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateMatchTargetSequenceRequest{ConfigID: 43253, ConfigVersion: 7, Type: "website", TargetSequence: []struct {
				TargetID int "json:\"targetId\""
				Sequence int "json:\"sequence\""
			}{struct {
				TargetID int "json:\"targetId\""
				Sequence int "json:\"sequence\""
			}{TargetID: 2052813, Sequence: 1}, struct {
				TargetID int "json:\"targetId\""
				Sequence int "json:\"sequence\""
			}{TargetID: 2971336, Sequence: 2}}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResMatchTargetSequence/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target_sequence.test", "id", "43253:7"),
						),
						//ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResMatchTargetSequence/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target_sequence.test", "id", "43253:7"),
						),
						//ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
