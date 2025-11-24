package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataContracts(t *testing.T) {
	t.Parallel()
	t.Run("list contracts", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		ctrs := papi.ContractsItems{Items: []*papi.Contract{
			{
				ContractID:       "ctr_test1",
				ContractTypeName: "ctr_typ_name_test1",
			},
			{
				ContractID:       "ctr_test2",
				ContractTypeName: "ctr_typ_name_test2",
			},
		}}

		client.PAPI.On("GetContracts",
			testutils.MockContext,
		).Return(&papi.GetContractsResponse{Contracts: ctrs, AccountID: "act_test"}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDataContracts/contracts.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_contracts.akacontracts", "id", "act_test"),
					resource.TestCheckOutput("aka_contract_id1", "ctr_test1"),
					resource.TestCheckOutput("aka_contract_id2", "ctr_test2"),
					resource.TestCheckOutput("aka_contract_typ_name1", "ctr_typ_name_test1"),
					resource.TestCheckOutput("aka_contract_typ_name2", "ctr_typ_name_test2"),
				),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})
}
