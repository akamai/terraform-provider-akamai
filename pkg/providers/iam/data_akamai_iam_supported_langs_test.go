package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSupportedLangs(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.IAM.Test(testutils.TattleT{T: t})
		client.IAM.On("SupportedLanguages", testutils.MockContext).Return([]string{"first", "second", "third"}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_supported_langs.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.#", "3"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.0", "first"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.1", "second"),
						resource.TestCheckResourceAttr("data.akamai_iam_supported_langs.test", "languages.2", "third"),
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
		client.IAM.On("SupportedLanguages", testutils.MockContext).Return([]string{}, errors.New("Could not set supported languages in state"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`Could not set supported languages in state`),
				},
			},
		})

		client.IAM.AssertExpectations(t)
	})

}
