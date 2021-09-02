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
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetMatchTargetResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTarget.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crmt := appsec.CreateMatchTargetResponse{}
		expectJSMT := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTargetCreated.json"))
		json.Unmarshal([]byte(expectJSMT), &crmt)

		rmmt := appsec.RemoveMatchTargetResponse{}
		expectJSRMT := compactJSON(loadFixtureBytes("testdata/TestResMatchTarget/MatchTargetCreated.json"))
		json.Unmarshal([]byte(expectJSRMT), &rmmt)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&cr, nil)

		createMatchTargetJSON := loadFixtureBytes("testdata/TestResMatchTarget/CreateMatchTarget.json")
		client.On("CreateMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateMatchTargetRequest{Type: "", ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createMatchTargetJSON},
		).Return(&crmt, nil)

		updateMatchTargetJSON := loadFixtureBytes("testdata/TestResMatchTarget/UpdateMatchTarget.json")
		client.On("UpdateMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967, JsonPayloadRaw: updateMatchTargetJSON},
		).Return(&cu, nil)

		client.On("RemoveMatchTarget",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&rmmt, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResMatchTarget/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "43253:3008967"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResMatchTarget/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "43253:3008967"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
