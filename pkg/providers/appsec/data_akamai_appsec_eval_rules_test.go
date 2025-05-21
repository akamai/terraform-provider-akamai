package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalRules_data_basic(t *testing.T) {
	t.Run("match by Rules ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getEvalRulesResponse := appsec.GetEvalRulesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSEvalRules/EvalRules.json"), &getEvalRulesResponse)
		require.NoError(t, err)

		configs := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configs)
		require.NoError(t, err)

		client.On("GetEvalRules",
			testutils.MockContext,
			appsec.GetEvalRulesRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getEvalRulesResponse, nil)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configs, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSEvalRules/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_rules.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
