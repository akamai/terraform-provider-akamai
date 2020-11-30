package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalRuleActions_data_basic(t *testing.T) {
	t.Run("match by EvalRuleActions ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetEvalRuleActionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSEvalRuleActions/EvalRuleActions.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetEvalRuleActions",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalRuleActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSEvalRuleActions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_rule_actions.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
