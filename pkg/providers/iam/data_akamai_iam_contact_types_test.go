package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataContactTypes(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.IAM.Test(testutils.TattleT{T: t})
		client.IAM.On("SupportedContactTypes", testutils.MockContext).Return([]string{"first", "second", "third"}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_contact_types.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_contact_types.test", "contact_types.#", "3"),
						resource.TestCheckResourceAttr("data.akamai_iam_contact_types.test", "contact_types.0", "first"),
						resource.TestCheckResourceAttr("data.akamai_iam_contact_types.test", "contact_types.1", "second"),
						resource.TestCheckResourceAttr("data.akamai_iam_contact_types.test", "contact_types.2", "third"),
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
		client.IAM.On("SupportedContactTypes", testutils.MockContext).Return(nil, errors.New("failed to get supported contact types"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/%v/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`failed to get supported contact types`),
				},
			},
		})

		client.IAM.AssertExpectations(t)
	})
}
