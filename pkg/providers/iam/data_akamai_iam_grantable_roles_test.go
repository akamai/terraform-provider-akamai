package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestGrantableRoles(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})
		client.On("ListGrantableRoles", mock.Anything).Return([]iam.RoleGrantedRole{
			{Description: "A", RoleID: 1, RoleName: "Can print A"},
			{Description: "B", RoleID: 2, RoleName: "Can print B"},
		}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				IsUnitTest:        true,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet("data.akamai_iam_grantable_roles.test", "id"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.#", "2"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.0.granted_role_id", "1"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.0.name", "Can print A"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.0.description", "A"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.1.granted_role_id", "2"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.1.name", "Can print B"),
							resource.TestCheckResourceAttr("data.akamai_iam_grantable_roles.test", "grantable_roles.1.description", "B"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})
		client.On("ListGrantableRoles", mock.Anything).Return([]iam.RoleGrantedRole{}, errors.New("could not get grantable roles"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				IsUnitTest:        true,
				Steps: []resource.TestStep{
					{
						Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
						ExpectError: regexp.MustCompile(`could not get grantable roles`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
