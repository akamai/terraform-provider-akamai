package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGrantableRoles(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.IAM.Test(testutils.TattleT{T: t})
		client.IAM.On("ListGrantableRoles", testutils.MockContext).Return([]iam.RoleGrantedRole{
			{Description: "A", RoleID: 1, RoleName: "Can print A"},
			{Description: "B", RoleID: 2, RoleName: "Can print B"},
		}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
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

		client.IAM.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.IAM.Test(testutils.TattleT{T: t})
		client.IAM.On("ListGrantableRoles", testutils.MockContext).Return([]iam.RoleGrantedRole{}, errors.New("could not get grantable roles"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`could not get grantable roles`),
				},
			},
		})

		client.IAM.AssertExpectations(t)
	})
}
