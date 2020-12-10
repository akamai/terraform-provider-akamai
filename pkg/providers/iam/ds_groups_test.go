package iam

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDSGroups(t *testing.T) {
	t.Parallel()

	t.Run("groups can nest 50 levels deep", func(t *testing.T) {
		t.Parallel()
		prov := provider{}

		assert.Equal(t, 50, GroupsNestingDepth(prov.dsGroups()), "incorrect nesting depth")
	})

	t.Run("happy path/without actions", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})

		{
			req := iam.ListGroupsRequest{}

			group1 := MakeGroup("test group 1", 101, 100, nil, nil)
			group4 := MakeGroup("test group 4", 104, 103, nil, nil)
			group5 := MakeGroup("test group 5", 105, 102, nil, nil)
			group3 := MakeGroup("test group 3", 103, 102, []iam.Group{group4}, nil)
			group2 := MakeGroup("test group 2", 102, 100, []iam.Group{group3, group5}, nil)
			res := []iam.Group{group1, group2, group3}

			client.On("ListGroups", mock.Anything, req).Return(res, nil)
		}

		expectedG1 := map[string]string{
			"name":            "test group 1",
			"group_id":        "101",
			"parent_group_id": "100",
			"time_created":    "2020-01-01T00:00:00Z",
			"time_modified":   "2020-01-01T00:00:00Z",
			"modified_by":     "modifier@akamai.net",
			"created_by":      "creator@akamai.net",
			"delete_allowed":  "false",
			"edit_allowed":    "false",
		}
		expectedG2 := map[string]string{
			"name":            "test group 2",
			"group_id":        "102",
			"parent_group_id": "100",
			"time_created":    "2020-01-01T00:00:00Z",
			"time_modified":   "2020-01-01T00:00:00Z",
			"modified_by":     "modifier@akamai.net",
			"created_by":      "creator@akamai.net",
			"delete_allowed":  "false",
			"edit_allowed":    "false",
		}
		expectedG3 := map[string]string{
			"name":            "test group 3",
			"group_id":        "103",
			"parent_group_id": "102",
			"time_created":    "2020-01-01T00:00:00Z",
			"time_modified":   "2020-01-01T00:00:00Z",
			"modified_by":     "modifier@akamai.net",
			"created_by":      "creator@akamai.net",
			"delete_allowed":  "false",
			"edit_allowed":    "false",
		}

		// See note below
		// expectedG4 := map[string]string{
		// 	"name":            "test group 4",
		// 	"group_id":        "104",
		// 	"parent_group_id": "103",
		// 	"time_created":    "2020-01-01T00:00:00Z",
		// 	"time_modified":   "2020-01-01T00:00:00Z",
		// 	"modified_by":     "modifier@akamai.net",
		// 	"created_by":      "creator@akamai.net",
		// 	"delete_allowed":  "false",
		// 	"edit_allowed":    "false",
		// }
		// expectedG5 := map[string]string{
		// 	"name":            "test group 5",
		// 	"group_id":        "105",
		// 	"parent_group_id": "102",
		// 	"time_created":    "2020-01-01T00:00:00Z",
		// 	"time_modified":   "2020-01-01T00:00:00Z",
		// 	"modified_by":     "modifier@akamai.net",
		// 	"created_by":      "creator@akamai.net",
		// 	"delete_allowed":  "false",
		// 	"edit_allowed":    "false",
		// }

		p := provider{}
		p.SetCache(metaCache{})
		p.SetClient(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_iam_groups.test", "id", "akamai_iam_groups"),

						// First level groups
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*", expectedG1),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*", expectedG2),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*", expectedG3),

						// About the checks below: There is a bug in TestCheckTypeSetElemNestedAttrs() that causes a panic
						// PR submitted https://github.com/hashicorp/terraform-plugin-sdk/pull/648

						// Second level groups (group 3 is found at both first and second levels)
						// resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*.sub_groups.*", expectedG3),
						// resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*.sub_groups.*", expectedG5),

						// Third level groups
						// resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*.sub_groups.*.sub_groups.*", expectedG4),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})

	t.Run("happy path/with actions", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})

		{
			req := iam.ListGroupsRequest{Actions: true}
			actions0 := iam.GroupActions{Delete: false, Edit: false}
			actions1 := iam.GroupActions{Delete: false, Edit: true}
			actions2 := iam.GroupActions{Delete: true, Edit: false}
			actions3 := iam.GroupActions{Delete: true, Edit: true}

			group1 := MakeGroup("test group 1", 101, 100, nil, &actions0)
			group4 := MakeGroup("test group 4", 104, 103, nil, &actions1)
			group5 := MakeGroup("test group 5", 105, 102, nil, &actions2)
			group3 := MakeGroup("test group 3", 103, 102, []iam.Group{group4}, &actions3)
			group2 := MakeGroup("test group 2", 102, 100, []iam.Group{group3, group5}, &actions0)
			// res := []iam.Group{group1, group2, group3, group4, group5}
			res := []iam.Group{group1, group2, group3}

			client.On("ListGroups", mock.Anything, req).Return(res, nil)
		}

		expectedG1 := map[string]string{
			"name":            "test group 1",
			"group_id":        "101",
			"parent_group_id": "100",
			"time_created":    "2020-01-01T00:00:00Z",
			"time_modified":   "2020-01-01T00:00:00Z",
			"modified_by":     "modifier@akamai.net",
			"created_by":      "creator@akamai.net",
			"delete_allowed":  "false",
			"edit_allowed":    "false",
		}
		expectedG2 := map[string]string{
			"name":            "test group 2",
			"group_id":        "102",
			"parent_group_id": "100",
			"time_created":    "2020-01-01T00:00:00Z",
			"time_modified":   "2020-01-01T00:00:00Z",
			"modified_by":     "modifier@akamai.net",
			"created_by":      "creator@akamai.net",
			"delete_allowed":  "false",
			"edit_allowed":    "false",
		}
		expectedG3 := map[string]string{
			"name":            "test group 3",
			"group_id":        "103",
			"parent_group_id": "102",
			"time_created":    "2020-01-01T00:00:00Z",
			"time_modified":   "2020-01-01T00:00:00Z",
			"modified_by":     "modifier@akamai.net",
			"created_by":      "creator@akamai.net",
			"delete_allowed":  "true",
			"edit_allowed":    "true",
		}

		// See the note below
		// expectedG4 := map[string]string{
		// 	"name":            "test group 4",
		// 	"group_id":        "104",
		// 	"parent_group_id": "103",
		// 	"time_created":    "2020-01-01T00:00:00Z",
		// 	"time_modified":   "2020-01-01T00:00:00Z",
		// 	"modified_by":     "modifier@akamai.net",
		// 	"created_by":      "creator@akamai.net",
		// 	"delete_allowed":  "false",
		// 	"edit_allowed":    "true",
		// }
		// expectedG5 := map[string]string{
		// 	"name":            "test group 5",
		// 	"group_id":        "105",
		// 	"parent_group_id": "102",
		// 	"time_created":    "2020-01-01T00:00:00Z",
		// 	"time_modified":   "2020-01-01T00:00:00Z",
		// 	"modified_by":     "modifier@akamai.net",
		// 	"created_by":      "creator@akamai.net",
		// 	"delete_allowed":  "true",
		// 	"edit_allowed":    "false",
		// }

		p := provider{}
		p.SetCache(metaCache{})
		p.SetClient(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_groups.test", "id"),

						// First level groups
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*", expectedG1),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*", expectedG2),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*", expectedG3),

						// NOTE: There is a bug in TestCheckTypeSetElemNestedAttrs() that causes a panic
						// PR submitted https://github.com/hashicorp/terraform-plugin-sdk/pull/648

						// Second level groups (group 3 is found at both first and second levels)
						// resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*.sub_groups.*", expectedG3),
						// resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*.sub_groups.*", expectedG5),

						// Third level groups
						// resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_groups.test", "groups.*.sub_groups.*.sub_groups.*", expectedG4),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
}

// counts the nesting depth of the groups in the groups resource schema
func GroupsNestingDepth(res *schema.Resource) int {
	if res, ok := res.Schema["sub_groups"]; ok {
		return 1 + GroupsNestingDepth(res.Elem.(*schema.Resource)) // Always a *schema.Resource for "sub_groups"
	}

	if res, ok := res.Schema["groups"]; ok {
		return 1 + GroupsNestingDepth(res.Elem.(*schema.Resource)) // Always a *schema.Resource for "groups"
	}

	return 0
}

// Convenience method to make a group
func MakeGroup(Name string, GroupID, PGroupID int64, SubGroups []iam.Group, Actions *iam.GroupActions) iam.Group {
	return iam.Group{
		Actions:       Actions,
		GroupName:     Name,
		GroupID:       GroupID,
		ParentGroupID: PGroupID,
		CreatedBy:     "creator@akamai.net",
		CreatedDate:   "2020-01-01T00:00:00Z",
		ModifiedBy:    "modifier@akamai.net",
		ModifiedDate:  "2020-01-01T00:00:00Z",
		SubGroups:     SubGroups,
	}
}
