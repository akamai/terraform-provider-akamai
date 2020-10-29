package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomRules_data_basic(t *testing.T) {
	t.Run("match by CustomRules ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetCustomRulesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSCustomRules/CustomRules.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetCustomRules",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomRulesRequest{ConfigID: 43253},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSCustomRules/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_custom_rules.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
