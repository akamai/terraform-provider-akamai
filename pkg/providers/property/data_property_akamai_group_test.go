package property

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccDataSourceGroup_basic(t *testing.T) {
	dataSourceName := "data.akamai_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
			{
				Config:      testAccDataSourceGroupNoContractWithGroupProvided(),
				ExpectError: regexp.MustCompile("^.+looking up group with name:.+contract ID is required for non-default name: .+$"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func TestVerifyDataSourceSchema(t *testing.T) {
	t.Run("test datasource ConflictsWith", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			Providers:  testAccProviders,
			IsUnitTest: true,
			Steps: []resource.TestStep{{
				Config:             testAccDataSourceContractConflicts(),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("\"contract\": conflicts with contract_id"),
			}},
		})
	})
}

func testAccDataSourceGroupBasic() string {
	return `
		provider "akamai" {
			papi_section = "papi"
			edgerc = "~/.edgerc"
		}

		data "akamai_group" "test" {
		}

		output "groupid" {
			value = "${data.akamai_group.test.id}"
		}
`
}

func testAccDataSourceContractConflicts() string {
	return `
		provider "akamai" {
			papi_section = "papi"
			edgerc = "~/.edgerc"
		}

		data "akamai_group" "test" {
            contract = "ctr_contract"
            contract_id = "ctr_contractId"
		}

		output "groupid" {
			value = "${data.akamai_group.test.id}"
		}
`
}

func testAccDataSourceGroupNoContractWithGroupProvided() string {
	return `
data "akamai_group" "test" {
	name = "Akamai Internal-3-984F"
}

output "groupid" {
value = "${data.akamai_group.test.id}"
}
`
}

func testAccCheckAkamaiGroupDestroy(_ *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Group Delete skipped ")

	return nil
}

func TestFindGroupByName(t *testing.T) {
	tests := map[string]struct {
		givenName              string
		givenContract          string
		givenGroups            []*papi.Group
		givenGroupsAccountName string
		isDefault              bool
		expected               *papi.Group
		withError              error
	}{
		"with default and no contract provided, return first group": {
			givenName:     "any name",
			givenContract: "",
			givenGroups: []*papi.Group{
				{
					GroupName: "Group A",
					GroupID:   "A",
				},
				{
					GroupName: "Group B",
					GroupID:   "B",
				},
			},
			givenGroupsAccountName: "",
			isDefault:              true,
			expected: &papi.Group{
				GroupName: "Group A",
				GroupID:   "A",
			},
		},
		"with default and no contract provided, no groups exist": {
			givenName:              "any name",
			givenContract:          "",
			givenGroups:            nil,
			givenGroupsAccountName: "",
			isDefault:              true,
			withError:              ErrNoGroupsFound,
		},
		"with default and contract provided, return matching group": {
			givenName:     "any name",
			givenContract: "ctr_B",
			givenGroups: []*papi.Group{
				{
					GroupName: "Group A",
					GroupID:   "A",
				},
				{
					GroupName: "Group B",
					GroupID:   "Account1-B",
				},
			},
			givenGroupsAccountName: "Account1",
			isDefault:              true,
			expected: &papi.Group{
				GroupName: "Group B",
				GroupID:   "Account1-B",
			},
		},
		"with default and contract provided, could not find group": {
			givenName:     "any name",
			givenContract: "ctr_B",
			givenGroups: []*papi.Group{
				{
					GroupName: "Group A",
					GroupID:   "A",
				},
			},
			givenGroupsAccountName: "Account1",
			isDefault:              true,
			withError:              ErrLookingUpGroupByName,
		},
		"not default and no contract provided, return error": {
			givenName:     "any name",
			givenContract: "",
			givenGroups: []*papi.Group{
				{
					GroupName: "Group A",
					GroupID:   "A",
				},
			},
			givenGroupsAccountName: "",
			isDefault:              false,
			withError:              ErrNoContractProvided,
		},
		"not default and contract provided, return matching group": {
			givenName:     "Group A",
			givenContract: "ctr_1",
			givenGroups: []*papi.Group{
				{
					GroupName: "Group A",
					GroupID:   "A",
				},
				{
					GroupName:   "Group A",
					GroupID:     "B",
					ContractIDs: []string{"ctr_1"},
				},
			},
			givenGroupsAccountName: "",
			isDefault:              false,
			expected: &papi.Group{
				GroupName:   "Group A",
				GroupID:     "B",
				ContractIDs: []string{"ctr_1"},
			},
		},
		"not default and contract provided, group does not belong to contract": {
			givenName:     "Group A",
			givenContract: "ctr_1",
			givenGroups: []*papi.Group{
				{
					GroupName: "Group A",
					GroupID:   "A",
				},
			},
			givenGroupsAccountName: "",
			isDefault:              false,
			withError:              ErrGroupNotInContract,
		},
		"not default and contract provided, no groups found by name": {
			givenName:              "Group A",
			givenContract:          "ctr_1",
			givenGroups:            []*papi.Group{},
			givenGroupsAccountName: "",
			isDefault:              false,
			withError:              ErrGroupNotInContract,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			groups := &papi.GetGroupsResponse{
				Groups: papi.GroupItems{
					Items: make([]*papi.Group, 0),
				},
			}
			groups.AccountName = test.givenGroupsAccountName
			for _, group := range test.givenGroups {
				groups.Groups.Items = append(groups.Groups.Items, group)
			}
			res, err := findGroupByName(test.givenName, test.givenContract, groups, test.isDefault)
			if test.withError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, res)
			require.NoError(t, err)
			assert.Equal(t, test.expected.GroupID, res.GroupID)
		})
	}
}

func TestDisabledCache(t *testing.T) {
	t.Run("testing cache_enabled=false", func(t *testing.T) {
		client := &mockpapi{}
		client.On("GetGroups")
		client.On("GetGroups", AnyCTX).Return(&papi.GetGroupsResponse{
			AccountID: "act_1-1TJZFB", AccountName: "example.com",
			Groups: papi.GroupItems{Items: []*papi.Group{{
				GroupID:       "grp_Example.com-1-1TJZH5",
				GroupName:     "Example.com-1-1TJZH5",
				ParentGroupID: "grp_parent",
				ContractIDs:   []string{"ctr_1234"},
			}}}}, nil)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDisabledCache/datasource-nocache.tf"),
				}},
			})
		})
	})
}
