package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomRules_data_basic(t *testing.T) {
	t.Run("match by CustomRules ID", func(t *testing.T) {
		client := &mockappsec{}

		getCustomRulesResponse := appsec.GetCustomRulesResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestDSCustomRules/CustomRules.json"), &getCustomRulesResponse)
		require.NoError(t, err)

		client.On("GetCustomRules",
			mock.Anything,
			appsec.GetCustomRulesRequest{ConfigID: 43253},
		).Return(&getCustomRulesResponse, nil)

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
