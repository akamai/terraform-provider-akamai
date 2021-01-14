package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mock "github.com/stretchr/testify/mock"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
)

func TestDSSupportedLangs(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("SupportedLanguages", mock.Anything).Return([]string{"first", "second", "third"}, nil)

		p := provider{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_supported_langs.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.0", "first"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.1", "second"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.2", "third"),
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
		client.On("SupportedLanguages", mock.Anything).Return([]string{}, errors.New("Could not set supported languages in state"))

		p := provider{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`Could not set supported languages in state`),
				},
			},
		})

		client.AssertExpectations(t)
	})

}
