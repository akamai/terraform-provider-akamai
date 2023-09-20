package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiContractsGroups_data_basic(t *testing.T) {
	t.Run("match by ContractsGroups ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getContractsGroupsResponse := appsec.GetContractsGroupsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSContractsGroups/ContractsGroups.json"), &getContractsGroupsResponse)
		require.NoError(t, err)

		client.On("GetContractsGroups",
			mock.Anything,
			appsec.GetContractsGroupsRequest{},
		).Return(&getContractsGroupsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSContractsGroups/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_contracts_groups.test", "id", "C-1FRYVV3"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
