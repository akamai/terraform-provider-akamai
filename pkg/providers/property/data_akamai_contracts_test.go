package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataContracts(t *testing.T) {
	t.Run("list contracts", func(t *testing.T) {
		client := &papi.Mock{}
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

		client.On("GetContracts",
			mock.Anything,
		).Return(&papi.GetContractsResponse{Contracts: ctrs, AccountID: "act_test"}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataContracts/contracts.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contracts.akacontracts", "id", "act_test"),
						resource.TestCheckOutput("aka_contract_id1", "ctr_test1"),
						resource.TestCheckOutput("aka_contract_id2", "ctr_test2"),
						resource.TestCheckOutput("aka_contract_typ_name1", "ctr_typ_name_test1"),
						resource.TestCheckOutput("aka_contract_typ_name2", "ctr_typ_name_test2"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})
}
