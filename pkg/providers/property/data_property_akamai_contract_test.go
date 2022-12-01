package property

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
)

func Test_DSReadContract(t *testing.T) {
	t.Run("read contract with group_id in group provided", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_id_in_group.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
					),
				}},
			})
		})
	})

	t.Run("read contract with group id w/o prefix in group", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_id_in_group_wo_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
					),
				}},
			})
		})
	})

	t.Run("read contract with group name in group", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_name_in_group.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
					),
				}},
			})
		})
	})

	t.Run("read contract with group name and group conflict", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_name_and_group.tf"),
					ExpectError: regexp.MustCompile("only one of `group,group_id,group_name` can be specified"),
				}},
			})
		})
	})

	t.Run("read contract with group id provided", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
					),
				}},
			})
		})
	})

	t.Run("read contract with group id without prefix", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_id_without_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
					),
				}},
			})
		})
	})

	t.Run("read contract with group name", func(t *testing.T) {
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
				{
					GroupID:       "grp_12346",
					GroupName:     "default",
					ParentGroupID: "grp_parent",
					ContractIDs:   []string{"ctr_1234"},
				},
			}}}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSContractRequired/ds_contract_with_group_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
						resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
					),
				}},
			})
		})
	})
}
