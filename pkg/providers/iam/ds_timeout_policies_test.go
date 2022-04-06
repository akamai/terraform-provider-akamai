package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mock "github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
)

func TestDSTimeoutPolicies(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})

		res := []iam.TimeoutPolicy{
			{Name: "first", Value: 11},
			{Name: "second", Value: 22},
			{Name: "third", Value: 33},
		}
		client.On("ListTimeoutPolicies", mock.Anything).Return(res, nil)

		p := providerOld{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_timeout_policies.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.first", "11"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.second", "22"),
						resource.TestCheckResourceAttr("data.akamai_iam_timeout_policies.test", "policies.third", "33"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("ListTimeoutPolicies", mock.Anything).Return(nil, errors.New("Could not get supported timeout policies"))

		p := providerOld{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`Could not get supported timeout policies`),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
