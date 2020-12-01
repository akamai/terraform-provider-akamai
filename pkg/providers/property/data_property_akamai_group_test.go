package property

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func Test_DSReadGroup(t *testing.T) {
	t.Run("read group with name and contract provided", func(t *testing.T) {
		client := &mockpapi{}
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
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSGroup/ds-group-w-name-and-contract.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "group_name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract_id", "ctr_1234"),
					),
				}},
			})
		})
	})

	t.Run("read group with name and contract_id provided", func(t *testing.T) {
		client := &mockpapi{}
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
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSGroup/ds-group-w-name-and-contract_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "group_name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract_id", "ctr_1234"),
					),
				}},
			})
		})
	})

	t.Run("read group with group_name and contract_id provided", func(t *testing.T) {
		client := &mockpapi{}
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
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "group_name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract_id", "ctr_1234"),
					),
				}},
			})
		})
	})

	t.Run("read group with group_name and contract provided", func(t *testing.T) {
		client := &mockpapi{}
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
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSGroup/ds-group-w-group-name-and-contract.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "group_name", "Example.com-1-1TJZH5"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract_id", "ctr_1234"),
					),
				}},
			})
		})
	})

	t.Run("read group with group_name and name conflict", func(t *testing.T) {
		client := &mockpapi{}
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
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSGroup/ds-group-w-group-name-and-name-conflict.tf"),
					ExpectError: regexp.MustCompile("only one of `group_name,name` can be specified"),
				}},
			})
		})
	})

	t.Run("read group with contract_id and contract conflict", func(t *testing.T) {
		client := &mockpapi{}
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
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSGroup/ds-group-w-contract-id-and-contract-conflict.tf"),
					ExpectError: regexp.MustCompile("only one of `contract,contract_id` can be specified"),
				}},
			})
		})
	})
}
