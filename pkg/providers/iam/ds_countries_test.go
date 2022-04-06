package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mock "github.com/stretchr/testify/mock"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
)

func TestDSCountries(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("SupportedCountries", mock.Anything).Return([]string{"first", "second", "third"}, nil)

		p := providerOld{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_countries.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_countries.test", "countries.0", "first"),
						resource.TestCheckResourceAttr("data.akamai_iam_countries.test", "countries.1", "second"),
						resource.TestCheckResourceAttr("data.akamai_iam_countries.test", "countries.2", "third"),
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
		client.On("SupportedCountries", mock.Anything).Return([]string{}, errors.New("Could not get supported countries"))

		p := providerOld{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`Could not get supported countries`),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
