package iam

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestGroupDataSource(t *testing.T) {
	mockGetGroup := func(client *iam.Mock, group *iam.Group, times int) *mock.Call {
		return client.On("GetGroup", testutils.MockContext, iam.GetGroupRequest{
			GroupID: 123,
			Actions: true,
		}).Return(group, nil).Times(times)
	}

	tests := map[string]struct {
		init           func(mock *iam.Mock)
		expectedError  *regexp.Regexp
		expectedChecks resource.TestCheckFunc
		givenTF        string
	}{
		"normal case - all fields": {
			init: func(client *iam.Mock) {
				group := &iam.Group{
					GroupID:       123,
					GroupName:     "parent_group",
					CreatedBy:     "DevUser",
					CreatedDate:   test.NewTimeFromString(t, "2024-05-28T06:58:26Z"),
					ModifiedBy:    "TestUser",
					ModifiedDate:  test.NewTimeFromString(t, "2024-05-28T06:58:27Z"),
					ParentGroupID: 0,
					Actions: &iam.GroupActions{
						Delete: true,
						Edit:   true,
					},
					SubGroups: []iam.Group{},
				}
				mockGetGroup(client, group, 3)
			},
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "group_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "group_name", "parent_group"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "created_by", "DevUser"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "modified_by", "TestUser"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "parent_group_id", "0"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "actions.delete", "true"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "actions.edit", "true"),
			),
			givenTF: "valid.tf",
		},
		"normal case - nested subgroups": {
			init: func(client *iam.Mock) {
				group := &iam.Group{
					GroupID:       123,
					GroupName:     "parent_group",
					CreatedBy:     "DevUser",
					CreatedDate:   test.NewTimeFromString(t, "2024-05-28T06:58:26Z"),
					ModifiedBy:    "TestUser",
					ModifiedDate:  test.NewTimeFromString(t, "2024-05-28T06:58:27Z"),
					ParentGroupID: 0,
					SubGroups: []iam.Group{
						{
							GroupID:       456,
							GroupName:     "child_group",
							CreatedBy:     "creator_child",
							ModifiedBy:    "modifier_child",
							ParentGroupID: 123,
						},
					},
				}
				mockGetGroup(client, group, 3)
			},
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "group_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "group_name", "parent_group"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "sub_groups.0.group_id", "456"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "sub_groups.0.group_name", "child_group"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "sub_groups.0.created_by", "creator_child"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "sub_groups.0.modified_by", "modifier_child"),
				resource.TestCheckResourceAttr("data.akamai_iam_group.test", "sub_groups.0.parent_group_id", "123"),
			),
			givenTF: "valid.tf",
		},
		"error - too deep nesting": {
			init: func(client *iam.Mock) {
				group := &iam.Group{
					GroupID:       123,
					GroupName:     "root_group",
					CreatedBy:     "DevUser",
					ModifiedBy:    "TestUser",
					ParentGroupID: 0,
					SubGroups:     generateDeepSubGroups(51),
				}
				mockGetGroup(client, group, 1)
			},
			expectedError: regexp.MustCompile("unsupported subgroup depth"),
			givenTF:       "valid.tf",
		},
		"api failed": {
			init: func(client *iam.Mock) {
				client.On("GetGroup", testutils.MockContext, iam.GetGroupRequest{
					GroupID: 123,
					Actions: true,
				}).Return(nil, errors.New("api failed")).Once()
			},
			expectedError: regexp.MustCompile("api failed"),
			givenTF:       "valid.tf",
		},
		"missing group_id": {
			expectedError: regexp.MustCompile("The argument \"group_id\" is required, but no definition was found."),
			givenTF:       "missing_group_id.tf",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataGroup/%s", tc.givenTF),
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

func generateDeepSubGroups(depth int64) []iam.Group {
	if depth == 0 {
		return nil
	}

	return []iam.Group{
		{
			GroupID:       depth,
			GroupName:     fmt.Sprintf("group%d", depth),
			CreatedBy:     "DevUser",
			ModifiedBy:    "TestUser",
			ParentGroupID: depth - 1,
			SubGroups:     generateDeepSubGroups(depth - 1),
		},
	}
}
