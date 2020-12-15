package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiMatchTarget_res_basic(t *testing.T) {
	t.Run("match by MatchTarget ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateMatchTargetResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTargetUpdated.json"))
		//expectJSU := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTarget.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetMatchTargetResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTarget.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crmt := appsec.CreateMatchTargetResponse{}
		expectJSMT := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTargetCreated.json"))
		//expectJSMT := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTarget.json"))
		json.Unmarshal([]byte(expectJSMT), &crmt)

		rmmt := appsec.RemoveMatchTargetResponse{}
		expectJSRMT := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTargetCreated.json"))
		json.Unmarshal([]byte(expectJSRMT), &rmmt)

		client.On("GetMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 15, TargetID: 3008967},
		).Return(&cr, nil)

		client.On("CreateMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateMatchTargetRequest{Type: "website", ConfigID: 43253, ConfigVersion: 15, DefaultFile: "NO_MATCH", EffectiveSecurityControls: struct {
				ApplyApplicationLayerControls bool "json:\"applyApplicationLayerControls\""
				ApplyBotmanControls           bool "json:\"applyBotmanControls\""
				ApplyNetworkLayerControls     bool "json:\"applyNetworkLayerControls\""
				ApplyRateControls             bool "json:\"applyRateControls\""
				ApplyReputationControls       bool "json:\"applyReputationControls\""
				ApplySlowPostControls         bool "json:\"applySlowPostControls\""
			}{ApplyApplicationLayerControls: false, ApplyBotmanControls: false, ApplyNetworkLayerControls: false, ApplyRateControls: false, ApplyReputationControls: false, ApplySlowPostControls: false}, FileExtensions: []string{"carb", "pct", "pdf", "swf", "cct", "jpeg", "js", "wmls", "hdml", "pws"}, FilePaths: []string{"/cache/aaabbc*"}, Hostnames: []string{"m.example.com", "www.example.net", "example.com"}, IsNegativeFileExtensionMatch: false, IsNegativePathMatch: false, SecurityPolicy: struct {
				PolicyID string "json:\"policyId\""
			}{PolicyID: "AAAA_81230"}, Sequence: 1, BypassNetworkLists: []struct {
				Name string "json:\"name\""
				ID   string "json:\"id\""
			}(nil)},
		).Return(&crmt, nil)

		client.On("UpdateMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateMatchTargetRequest{Type: "website", ConfigID: 43253, ConfigVersion: 15, TargetID: 3008967, DefaultFile: "NO_MATCH", EffectiveSecurityControls: struct {
				ApplyApplicationLayerControls bool "json:\"applyApplicationLayerControls\""
				ApplyBotmanControls           bool "json:\"applyBotmanControls\""
				ApplyNetworkLayerControls     bool "json:\"applyNetworkLayerControls\""
				ApplyRateControls             bool "json:\"applyRateControls\""
				ApplyReputationControls       bool "json:\"applyReputationControls\""
				ApplySlowPostControls         bool "json:\"applySlowPostControls\""
			}{ApplyApplicationLayerControls: false, ApplyBotmanControls: false, ApplyNetworkLayerControls: false, ApplyRateControls: false, ApplyReputationControls: false, ApplySlowPostControls: false}, FileExtensions: []string{"carb", "pct", "pdf", "swf", "cct", "jpeg", "js", "wmls", "hdml", "pws"}, FilePaths: []string{"/cache/aaabbc*"}, Hostnames: []string{"m1.example.com", "www.example.net", "example.com"}, IsNegativeFileExtensionMatch: false, IsNegativePathMatch: false, SecurityPolicy: struct {
				PolicyID string "json:\"policyId\""
			}{PolicyID: "AAAA_81230"}, Sequence: 1, BypassNetworkLists: []struct {
				Name string "json:\"name\""
				ID   string "json:\"id\""
			}(nil)},
		).Return(&cu, nil)

		client.On("RemoveMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveMatchTargetRequest{ConfigID: 43253, ConfigVersion: 15, TargetID: 3008967},
		).Return(&rmmt, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResMatchTarget/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "3008967"), //3008967
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResMatchTarget/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "3008967"),
							//resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "is_negative_file_extension_match", "false"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
