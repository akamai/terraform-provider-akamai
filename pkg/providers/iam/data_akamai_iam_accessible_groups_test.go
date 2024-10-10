package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccessibleGroups(t *testing.T) {
	mockPositiveCase := func(client *iam.Mock, returnedGroups iam.ListAccessibleGroupsResponse, times int) *mock.Call {
		return client.On("ListAccessibleGroups", mock.Anything, iam.ListAccessibleGroupsRequest{UserName: "user1"}).Return(returnedGroups, nil).Times(times)
	}

	generateCheckForGroup := func(path, groupID, groupName, isBlocked, roleDescription, roleID, roleName, subGroupsNumber string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.group_id", path), groupID),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.group_name", path), groupName),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.is_blocked", path), isBlocked),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.role_description", path), roleDescription),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.role_id", path), roleID),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.role_name", path), roleName),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.sub_groups.#", path), subGroupsNumber),
		)
	}

	generateCheckForSubGroup := func(path, groupID, groupName, parentGroupID, subGroupsNumber string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.group_id", path), groupID),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.group_name", path), groupName),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.parent_group_id", path), parentGroupID),
			resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", fmt.Sprintf("%s.sub_groups.#", path), subGroupsNumber),
		)
	}

	generateSubGroups := func(maxDepth, parentGroupID int64) []iam.AccessibleSubGroup {
		var childSubgroups []iam.AccessibleSubGroup

		for groupID := maxDepth; groupID > 1; groupID-- {
			currentSubgroups := iam.AccessibleSubGroup{
				GroupID:       groupID,
				GroupName:     fmt.Sprintf("group%d", groupID),
				ParentGroupID: groupID - 1,
				SubGroups:     childSubgroups,
			}
			childSubgroups = []iam.AccessibleSubGroup{currentSubgroups}
		}
		return []iam.AccessibleSubGroup{
			{
				GroupID:       1,
				GroupName:     "group1",
				ParentGroupID: parentGroupID,
				SubGroups:     childSubgroups,
			},
		}
	}

	generateIntermediateChecksForGeneratedSubGroup := func(min, max int) resource.TestCheckFunc {
		var testCases []resource.TestCheckFunc

		for i := min; i < max; i++ {
			testCases = append(testCases, generateCheckForSubGroup(fmt.Sprintf("accessible_groups.0%s", strings.Repeat(".sub_groups.0", i)), strconv.Itoa(i), fmt.Sprintf("group%d", i), strconv.Itoa(i-1), "1"))
		}

		return resource.ComposeAggregateTestCheckFunc(testCases...)
	}

	tests := map[string]struct {
		init           func(mock *iam.Mock)
		config         string
		expectedError  *regexp.Regexp
		expectedChecks resource.TestCheckFunc
	}{
		"normal case - all fields": {
			init: func(client *iam.Mock) {
				returnedGroups := []iam.AccessibleGroup{
					{
						GroupID:         123,
						GroupName:       "first_group",
						RoleID:          456,
						RoleName:        "admin",
						RoleDescription: "admin description",
						IsBlocked:       true,
					},
					{
						GroupID:         321,
						GroupName:       "second_group",
						RoleID:          654,
						RoleName:        "group1",
						RoleDescription: "group1 description",
						IsBlocked:       false,
						SubGroups: []iam.AccessibleSubGroup{
							{
								GroupID:       234,
								GroupName:     "name2",
								ParentGroupID: 321,
								SubGroups: []iam.AccessibleSubGroup{
									{
										GroupID:       345,
										GroupName:     "name3",
										ParentGroupID: 234,
										SubGroups:     []iam.AccessibleSubGroup{},
									},
								},
							},
						},
					},
				}
				mockPositiveCase(client, returnedGroups, 3)
			},
			config: "testdata/TestDataAccessibleGroups/basic.tf",
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", "username", "user1"),
				resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", "accessible_groups.#", "2"),
				generateCheckForGroup("accessible_groups.0", "123", "first_group", "true", "admin description", "456", "admin", "0"),
				generateCheckForGroup("accessible_groups.1", "321", "second_group", "false", "group1 description", "654", "group1", "1"),
				generateCheckForSubGroup("accessible_groups.1.sub_groups.0", "234", "name2", "321", "1"),
				generateCheckForSubGroup("accessible_groups.1.sub_groups.0.sub_groups.0", "345", "name3", "234", "0"),
			),
		},
		"normal case - max depth": {
			init: func(client *iam.Mock) {
				returnedGroups := []iam.AccessibleGroup{
					{
						GroupID:         123,
						GroupName:       "root group",
						RoleID:          321,
						RoleName:        "role321",
						RoleDescription: "role321 description",
						IsBlocked:       false,
						SubGroups:       generateSubGroups(50, 123),
					},
				}
				mockPositiveCase(client, returnedGroups, 3)
			},
			config: "testdata/TestDataAccessibleGroups/basic.tf",
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", "username", "user1"),
				resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", "accessible_groups.#", "1"),
				generateCheckForSubGroup(fmt.Sprintf("accessible_groups.0%s", strings.Repeat(".sub_groups.0", 1)), "1", "group1", "123", "1"),
				generateIntermediateChecksForGeneratedSubGroup(2, 50),
				generateCheckForSubGroup(fmt.Sprintf("accessible_groups.0%s", strings.Repeat(".sub_groups.0", 50)), "50", "group50", "49", "0"),
			),
		},
		"normal case - empty response": {
			init: func(client *iam.Mock) {
				mockPositiveCase(client, []iam.AccessibleGroup{}, 3)
			},
			config: "testdata/TestDataAccessibleGroups/basic.tf",
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", "username", "user1"),
				resource.TestCheckResourceAttr("data.akamai_iam_accessible_groups.groups", "accessible_groups.#", "0"),
			),
		},
		"error - too deep nesting": {
			init: func(client *iam.Mock) {
				returnedGroups := []iam.AccessibleGroup{
					{
						GroupID:         123,
						GroupName:       "root group",
						RoleID:          321,
						RoleName:        "role321",
						RoleDescription: "role321 description",
						IsBlocked:       false,
						SubGroups:       generateSubGroups(51, 123),
					},
				}
				mockPositiveCase(client, returnedGroups, 1)
			},
			config:        "testdata/TestDataAccessibleGroups/basic.tf",
			expectedError: regexp.MustCompile("unsupported subgroup depth"),
		},
		"error - api failed": {
			init: func(client *iam.Mock) {
				client.On("ListAccessibleGroups", mock.Anything, iam.ListAccessibleGroupsRequest{UserName: "user1"}).Return(nil, errors.New("api failed")).Times(1)
			},
			config:        "testdata/TestDataAccessibleGroups/basic.tf",
			expectedError: regexp.MustCompile("api failed"),
		},
		"error - missing username": {
			init:          func(client *iam.Mock) {},
			config:        "testdata/TestDataAccessibleGroups/no-username.tf",
			expectedError: regexp.MustCompile("Missing required argument"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			tc.init(client)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.config),
							Check:       tc.expectedChecks,
							ExpectError: tc.expectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
