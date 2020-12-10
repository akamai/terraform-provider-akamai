package iam

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mock "github.com/stretchr/testify/mock"
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

		p := provider{}
		p.SetCache(metaCache{})
		p.SetClient(client)

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
}
