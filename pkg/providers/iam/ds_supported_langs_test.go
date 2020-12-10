package iam

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mock "github.com/stretchr/testify/mock"
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
						resource.TestCheckTypeSetElemAttr("data.akamai_iam_supported_langs.test", "languages.*", "first"),
						resource.TestCheckTypeSetElemAttr("data.akamai_iam_supported_langs.test", "languages.*", "second"),
						resource.TestCheckTypeSetElemAttr("data.akamai_iam_supported_langs.test", "languages.*", "third"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
