package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mock "github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
)

func TestDataTimeoutPolicies(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		client := &mockiam{}
		client.Test(test.TattleT{T: t})

		res := []iam.TimeoutPolicy{
			{Name: "first", Value: 11},
			{Name: "second", Value: 22},
			{Name: "third", Value: 33},
		}
		client.On("ListTimeoutPolicies", mock.Anything).Return(res, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
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
		client := &mockiam{}
		client.Test(test.TattleT{T: t})
		client.On("ListTimeoutPolicies", mock.Anything).Return(nil, errors.New("Could not get supported timeout policies"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{
					{
						Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
						ExpectError: regexp.MustCompile(`Could not get supported timeout policies`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
