package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomRules_data_basic(t *testing.T) {
	t.Run("match by CustomRules ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getCustomRulesResponse := appsec.GetCustomRulesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSCustomRules/CustomRules.json"), &getCustomRulesResponse)
		require.NoError(t, err)

		client.On("GetCustomRules",
			testutils.MockContext,
			appsec.GetCustomRulesRequest{ConfigID: 43253},
		).Return(&getCustomRulesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSCustomRules/match_by_id.tf"),
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
