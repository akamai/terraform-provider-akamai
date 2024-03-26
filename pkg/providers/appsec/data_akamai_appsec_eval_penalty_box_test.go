package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBox_data_basic(t *testing.T) {
	t.Run("match by EvalPenaltyBox ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		penaltyBox := appsec.GetPenaltyBoxResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSEvalPenaltyBox/PenaltyBox.json"), &penaltyBox)
		require.NoError(t, err)

		client.On("GetEvalPenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&penaltyBox, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSEvalPenaltyBox/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_penalty_box.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
