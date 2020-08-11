package akamai

import (
	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceGroup_basic(t *testing.T) {
	dataSourceName := "data.akamai_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceGroup_basic() string {
	return `
provider "akamai" {
  papi_section = "papi"
}

data "akamai_group" "test" {
}

output "groupid" {
value = "${data.akamai_group.test.id}"
}
`
}

func testAccCheckAkamaiGroupDestroy(s *terraform.State) error {
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
			withError:              ErrPapiNoGroupsFound,
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
			withError:              ErrPapiGroupNotFound,
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
			withError:              ErrPapiNoContractProvided,
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
			withError:              ErrPapiGroupNotInContract,
		},
		"not default and contract provided, no groups found by name": {
			givenName:              "Group A",
			givenContract:          "ctr_1",
			givenGroups:            []*papi.Group{},
			givenGroupsAccountName: "",
			isDefault:              false,
			withError:              ErrPapiFindingGroupsByName,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			groups := papi.NewGroups()
			groups.AccountName = test.givenGroupsAccountName
			for _, group := range test.givenGroups {
				groups.AddGroup(group)
			}
			res, err := findGroupByName(test.givenName, test.givenContract, groups, test.isDefault)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected.GroupID, res.GroupID)
		})
	}
}
