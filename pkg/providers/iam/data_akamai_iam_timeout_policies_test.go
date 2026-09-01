package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataTimeoutPolicies(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.IAM.Test(testutils.TattleT{T: t})

		res := []iam.TimeoutPolicy{
			{Name: "first", Value: 11},
			{Name: "second", Value: 22},
			{Name: "third", Value: 33},
		}
		client.IAM.On("ListTimeoutPolicies", testutils.MockContext).Return(res, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_timeout_policies.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.%", "3"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.first", "11"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.second", "22"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.third", "33"),
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
		client.IAM.On("ListTimeoutPolicies", testutils.MockContext).Return(nil, errors.New("Could not get supported timeout policies"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`Could not get supported timeout policies`),
				},
			},
		})

		client.IAM.AssertExpectations(t)
	})
}
