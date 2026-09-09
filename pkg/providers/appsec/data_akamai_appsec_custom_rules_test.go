package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomRules_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by CustomRules ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getCustomRulesResponse := appsec.GetCustomRulesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSCustomRules/CustomRules.json"), &getCustomRulesResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetCustomRules",
			testutils.MockContext,
			appsec.GetCustomRulesRequest{ConfigID: 43253},
		).Return(&getCustomRulesResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSCustomRules/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_custom_rules.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
