package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBox_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by EvalPenaltyBox ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		penaltyBox := appsec.GetPenaltyBoxResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSEvalPenaltyBox/PenaltyBox.json"), &penaltyBox)
		require.NoError(t, err)

		client.APPSEC.On("GetEvalPenaltyBox",
			testutils.MockContext,
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&penaltyBox, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSEvalPenaltyBox/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_eval_penalty_box.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
