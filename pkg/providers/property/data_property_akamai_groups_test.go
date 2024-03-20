package property

import (
	"log"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDataSourceMultipleGroups_basic(t *testing.T) {
	t.Run("test output", func(t *testing.T) {
		client := &papi.Mock{}
		contractIDs := []string{"ctr_1234"}
		groups := []map[string]interface{}{{
			"group_id":        "grp_12345",
			"group_name":      "Example.com-1-1TJZH5",
			"parent_group_id": "grp_parent",
			"contractIds":     contractIDs,
		}}
		client.On("GetGroups", AnyCTX).Return(&papi.GetGroupsResponse{
			AccountID: "act_1-1TJZFB", AccountName: "example.com",
			Groups: papi.GroupItems{Items: []*papi.Group{{
				GroupID:       groups[0]["group_id"].(string),
				GroupName:     groups[0]["group_name"].(string),
				ParentGroupID: groups[0]["parent_group_id"].(string),
				ContractIDs:   contractIDs,
			}}}}, nil)
		useClient(client, nil, func() {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				CheckDestroy:             testAccCheckAkamaiMultipleGroupsDestroy,
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config: testAccDataSourceMultipleGroupsBasic(),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("group_id1", groups[0]["group_id"].(string)),
							resource.TestCheckOutput("group_name1", groups[0]["group_name"].(string)),
							resource.TestCheckOutput("parent_group_id1", groups[0]["parent_group_id"].(string)),
							resource.TestCheckOutput("group_contract1", contractIDs[0]),
						),
					},
				},
			})
		})
	})
}

func TestGroup_ContractNotFoundInState(t *testing.T) {
	t.Run("contractId not found in state", func(t *testing.T) {
		client := &papi.Mock{}
		contractIDs := []string{"ctr_contractID"}
		groups := []map[string]interface{}{{
			"group_id":        "grp_test",
			"group_name":      "test",
			"parent_group_id": "grp_parent",
			"contractIds":     contractIDs,
		}}
		client.On("GetGroups", AnyCTX).Return(&papi.GetGroupsResponse{
			AccountID: "act_1-1TJZFB", AccountName: "example.com",
			Groups: papi.GroupItems{Items: []*papi.Group{{
				GroupID:       groups[0]["group_id"].(string),
				GroupName:     groups[0]["group_name"].(string),
				ParentGroupID: groups[0]["parent_group_id"].(string),
				ContractIDs:   contractIDs,
			}}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSContractRequired/groups.tf"),
				}},
			})
		})
	})
}

func testAccDataSourceMultipleGroupsBasic() string {
	return `
		provider "akamai" {
			edgerc = "../../common/testutils/edgerc"
		}

		data "akamai_groups" "test" {}

		output "group_id1" {
			value = "${data.akamai_groups.test.groups[0].group_id}"
		}

		output "group_name1" {
			value = "${data.akamai_groups.test.groups[0].group_name}"
		}

		output "parent_group_id1" {
			value = "${data.akamai_groups.test.groups[0].parent_group_id}"
		}

		output "group_contract1" {
			value = "${data.akamai_groups.test.groups[0].contract_ids[0]}"
		}
`
}

func testAccCheckAkamaiMultipleGroupsDestroy(_ *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Group Delete skipped ")

	return nil
}
