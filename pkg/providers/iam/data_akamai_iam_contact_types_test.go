package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataContactTypes(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})
		client.On("SupportedContactTypes", mock.Anything).Return([]string{"first", "second", "third"}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})
		client.On("SupportedContactTypes", mock.Anything).Return(nil, errors.New("failed to get supported contact types"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      test.Fixture("testdata/%v/step0.tf", t.Name()),
						ExpectError: regexp.MustCompile(`failed to get supported contact types`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
