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

	t.Run("happy path/no args", func(t *testing.T) {
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
		}}

		req := iam.ListRolesRequest{}

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("ListRoles", mock.Anything, req).Return(roles, nil)

		p := provider{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_roles.test", "id"),
						resource.TestCheckTypeSetElemNestedAttrs("data.akamai_iam_roles.test", "roles.*", map[string]string{
							"name":          "test role name",
							"role_id":       "100",
							"description":   "role description",
							"type":          string(iam.RoleTypeStandard),
							"time_created":  "2020-01-01T00:00:00Z",
							"time_modified": "2020-01-01T00:00:00Z",
							"created_by":    "creator@akamai.net",
							"modified_by":   "modifier@akamai.net",
						}),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
