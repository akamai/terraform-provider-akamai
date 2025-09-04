package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataTimeoutPolicies(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(testutils.TattleT{T: t})

		res := []iam.TimeoutPolicy{
			{Name: "first", Value: 11},
			{Name: "second", Value: 22},
			{Name: "third", Value: 33},
		}
		client.On("ListTimeoutPolicies", testutils.MockContext).Return(res, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(testutils.TattleT{T: t})
		client.On("ListTimeoutPolicies", testutils.MockContext).Return(nil, errors.New("Could not get supported timeout policies"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
						ExpectError: regexp.MustCompile(`Could not get supported timeout policies`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
