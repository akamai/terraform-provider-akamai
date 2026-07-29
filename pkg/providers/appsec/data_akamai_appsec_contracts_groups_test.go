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

func TestAkamaiContractsGroups_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by ContractsGroups ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getContractsGroupsResponse := appsec.GetContractsGroupsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSContractsGroups/ContractsGroups.json"), &getContractsGroupsResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetContractsGroups",
			testutils.MockContext,
			appsec.GetContractsGroupsRequest{},
		).Return(&getContractsGroupsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSContractsGroups/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_contracts_groups.test", "id", "C-1FRYVV3"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
