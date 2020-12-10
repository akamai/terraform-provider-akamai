package iam

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDSRoles(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		roles := []iam.Role{{
			RoleName:        "test role name",
			RoleID:          100,
			RoleDescription: "role description",
			RoleType:        iam.RoleTypeStandard,
			CreatedBy:       "creator@akamai.net",
			CreatedDate:     "2020-01-01T00:00:00Z",
			ModifiedBy:      "modifier@akamai.net",
			ModifiedDate:    "2020-01-01T00:00:00Z",
			GrantedRoles: []iam.RoleGrantedRole{{
				RoleName:    "granted test role name",
				RoleID:      101,
				Description: "granted role description",
			}},
			Users: []iam.RoleUser{{
				UIIdentityID:  "300",
				Email:         "user@akamai.net",
				FirstName:     "user's first name",
				LastName:      "user's last name",
				AccountID:     "200",
				LastLoginDate: "2020-01-01T00:00:00Z",
			}},
			Actions: &iam.RoleAction{
				Delete: true,
				Edit:   true,
			},
		}}

		gid := int64(300)
		req := iam.ListRolesRequest{
			GroupID:       &gid,
			Actions:       true,
			IgnoreContext: true,
			Users:         true,
		}

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("ListRoles", mock.Anything, req).Return(roles, nil)

		p := provider{}
		p.SetCache(metaCache{})
		p.SetClient(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_roles.test", "id"),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_roles.test", "roles.*", map[string]string{
							"name":           "test role name",
							"role_id":        "100",
							"description":    "role description",
							"type":           string(iam.RoleTypeStandard),
							"time_created":   "2020-01-01T00:00:00Z",
							"time_modified":  "2020-01-01T00:00:00Z",
							"created_by":     "creator@akamai.net",
							"modified_by":    "modifier@akamai.net",
							"edit_allowed":   "true",
							"delete_allowed": "true",
						}),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_roles.test", "roles.*.granted_roles.*", map[string]string{
							"name":        "granted test role name",
							"role_id":     "101",
							"description": "granted role description",
						}),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_roles.test", "roles.*.users.*", map[string]string{
							"user_id":    "300",
							"email":      "user@akamai.net",
							"first_name": "user's first name",
							"last_name":  "user's last name",
							"account_id": "200",
							"last_login": "2020-01-01T00:00:00Z",
						}),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
