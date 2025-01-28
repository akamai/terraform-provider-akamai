package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalGroups_data_basic(t *testing.T) {
	t.Run("match by Eval Attack Group ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configs := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configs)
		require.NoError(t, err)

		client.On("GetEvalGroups",
			testutils.MockContext,
			appsec.GetAttackGroupsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(nil, fmt.Errorf("GetEvalGroups failed"))

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
						Config: testutils.LoadFixtureString(t, "testdata/TestDSEvalGroups/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_groups.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`GetEvalGroups failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAkamaiEvalGroups_data_error_retrieving_eval_groups(t *testing.T) {
	t.Run("match by Eval Attack Group ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configs := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configs)
		require.NoError(t, err)

		client.On("GetEvalGroups",
			testutils.MockContext,
			appsec.GetAttackGroupsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(nil, fmt.Errorf("GetEvalGroups failed"))

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
						Config: testutils.LoadFixtureString(t, "testdata/TestDSEvalGroups/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_groups.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`GetEvalGroups failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
