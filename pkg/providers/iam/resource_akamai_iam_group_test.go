package iam

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceGroup(t *testing.T) {
	groupCreate := iam.Group{
		GroupID:       3,
		GroupName:     "test",
		ParentGroupID: 1,
		SubGroups: []iam.Group{
			{GroupID: 4},
			{GroupID: 5},
			{GroupID: 6},
		},
	}
	groupUpdate := iam.Group{
		GroupID:       groupCreate.GroupID,
		GroupName:     "another test",
		ParentGroupID: 7,
		SubGroups:     groupCreate.SubGroups,
	}

	tests := map[string]struct {
		init  func(*iam.Mock)
		steps []resource.TestStep
	}{
		"creation error": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResGroup/basic/basic.tf"),
					ExpectError: regexp.MustCompile("group creation error"),
				},
			},
			init: func(m *iam.Mock) {
				expectResourceIAMGroupCreate(m, groupCreate.ParentGroupID, groupCreate.GroupName, &groupCreate, fmt.Errorf("group creation error"))
			},
		},
		"group read error": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResGroup/basic/basic.tf"),
					ExpectError: regexp.MustCompile("group read error"),
				},
			},
			init: func(m *iam.Mock) {
				// step 1
				// create
				expectResourceIAMGroupCreate(m, groupCreate.ParentGroupID, groupCreate.GroupName, &groupCreate, nil)
				expectResourceIAMGroupRead(m, groupCreate.GroupID, &groupCreate, fmt.Errorf("group read error")).Once()

				// delete
				expectResourceIAMGroupDelete(m, int(groupCreate.GroupID), nil)
			},
		},
		"basic": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResGroup/basic/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_group.test", "parent_group_id", strconv.FormatInt(groupCreate.ParentGroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "id", strconv.FormatInt(groupCreate.GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "name", groupCreate.GroupName),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.#", strconv.Itoa(len(groupCreate.SubGroups))),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.0", strconv.FormatInt(groupCreate.SubGroups[0].GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.1", strconv.FormatInt(groupCreate.SubGroups[1].GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.2", strconv.FormatInt(groupCreate.SubGroups[2].GroupID, 10)),
					),
				},
			},
			init: func(m *iam.Mock) {
				// step 1
				// create
				expectResourceIAMGroupCreate(m, groupCreate.ParentGroupID, groupCreate.GroupName, &groupCreate, nil)
				expectResourceIAMGroupRead(m, groupCreate.GroupID, &groupCreate, nil).Once()

				// read
				expectResourceIAMGroupRead(m, groupCreate.GroupID, &groupCreate, nil).Once()

				// delete
				expectResourceIAMGroupDelete(m, int(groupCreate.GroupID), nil)
			},
		},
		"update": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResGroup/basic/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_group.test", "parent_group_id", strconv.FormatInt(groupCreate.ParentGroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "id", strconv.FormatInt(groupCreate.GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "name", groupCreate.GroupName),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.#", strconv.Itoa(len(groupCreate.SubGroups))),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.0", strconv.FormatInt(groupCreate.SubGroups[0].GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.1", strconv.FormatInt(groupCreate.SubGroups[1].GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.2", strconv.FormatInt(groupCreate.SubGroups[2].GroupID, 10)),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResGroup/update/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_group.test", "parent_group_id", strconv.FormatInt(groupUpdate.ParentGroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "id", strconv.FormatInt(groupUpdate.GroupID, 10)),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "name", groupUpdate.GroupName),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.#", strconv.Itoa(len(groupUpdate.SubGroups))),
						resource.TestCheckResourceAttr("akamai_iam_group.test", "sub_groups.0", strconv.FormatInt(groupUpdate.SubGroups[0].GroupID, 10)),
					),
				},
			},
			init: func(m *iam.Mock) {
				// step 1
				// create
				expectResourceIAMGroupCreate(m, groupCreate.ParentGroupID, groupCreate.GroupName, &groupCreate, nil)
				expectResourceIAMGroupRead(m, groupCreate.GroupID, &groupCreate, nil).Once()

				// read
				expectResourceIAMGroupRead(m, groupCreate.GroupID, &groupCreate, nil).Once()

				// step 2
				// refresh
				expectResourceIAMGroupRead(m, groupCreate.GroupID, &groupCreate, nil).Once()
				// update
				expectResourceIAMGroupUpdate(m, groupUpdate, groupUpdate.GroupID, nil, nil)
				expectResourceIAMGroupRead(m, groupUpdate.GroupID, &groupUpdate, nil).Twice()

				// delete
				expectResourceIAMGroupDelete(m, int(groupCreate.GroupID), nil)
			},
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
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectResourceIAMGroupUpdate(m *iam.Mock, group iam.Group, sourceGroupID int64, updateNameError, moveGroupError error) {
	onUpdateGroupName := m.On("UpdateGroupName", testutils.MockContext, iam.GroupRequest{GroupName: group.GroupName, GroupID: group.GroupID})
	if updateNameError != nil {
		onUpdateGroupName.Return(nil, updateNameError)
		return
	}
	onUpdateGroupName.Return(nil, nil)

	m.On("MoveGroup", testutils.MockContext, iam.MoveGroupRequest{DestinationGroupID: group.ParentGroupID, SourceGroupID: sourceGroupID}).Return(moveGroupError)
}

func expectResourceIAMGroupDelete(m *iam.Mock, groupID int, errRemoveGroup error) {
	m.On("RemoveGroup", testutils.MockContext, iam.RemoveGroupRequest{GroupID: int64(groupID)}).Return(errRemoveGroup)
}

func expectResourceIAMGroupRead(m *iam.Mock, groupID int64, group *iam.Group, errRead error) *mock.Call {
	onGet := m.On("GetGroup", testutils.MockContext, iam.GetGroupRequest{GroupID: groupID})
	if errRead != nil {
		return onGet.Return(nil, errRead)
	}
	return onGet.Return(group, nil)
}

func expectResourceIAMGroupCreate(m *iam.Mock, parentGroupID int64, createGroupName string, groupCreate *iam.Group, errCreate error) {
	m.On("CreateGroup", testutils.MockContext, iam.GroupRequest{
		GroupID: parentGroupID, GroupName: createGroupName,
	}).Return(groupCreate, errCreate).Once()
}
