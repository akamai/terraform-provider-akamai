package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_DSReadGroup(t *testing.T) {
	t.Run("read group with group_name and contract_id provided", func(t *testing.T) {
		client := &papi.Mock{}
		client.On("GetGroups")
		client.On("GetGroups", AnyCTX).Return(&papi.GetGroupsResponse{
			AccountID: "act_1-1TJZFB", AccountName: "example.com",
			Groups: papi.GroupItems{Items: []*papi.Group{
				{
					GroupID:       "grp_12345",
					GroupName:     "Example.com-1-1TJZH5",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				IsUnitTest:               true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "group_name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract_id", "ctr_1234"),
					),
				}},
			})
		})
	})
}
