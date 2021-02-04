package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationAnalysis_res_basic(t *testing.T) {
	t.Run("match by ReputationAnalysis ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateReputationAnalysisResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResReputationAnalysis/ReputationAnalysisUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetReputationAnalysisResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResReputationAnalysis/ReputationAnalysis.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		cd := appsec.RemoveReputationAnalysisResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResReputationAnalysis/ReputationAnalysisDelete.json"))
		json.Unmarshal([]byte(expectJSD), &cd)

		client.On("GetReputationAnalysis",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetReputationAnalysisRequest{ConfigID: 43253, Version: 12, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateReputationAnalysis",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateReputationAnalysisRequest{ConfigID: 43253, Version: 12, PolicyID: "AAAA_81230", ForwardToHTTPHeader: true, ForwardSharedIPToHTTPHeaderAndSIEM: true},
		).Return(&cu, nil)

		client.On("RemoveReputationAnalysis",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveReputationAnalysisRequest{ConfigID: 43253, Version: 12, PolicyID: "AAAA_81230", ForwardToHTTPHeader: false, ForwardSharedIPToHTTPHeaderAndSIEM: false},
		).Return(&cd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationAnalysis/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_analysis.test", "id", "43253:12:AAAA_81230"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResReputationAnalysis/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_analysis.test", "id", "43253:12:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
