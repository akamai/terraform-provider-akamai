package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/test"
)

func TestDataStates(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})

		req := iam.ListStatesRequest{Country: "test country"}
		client.On("ListStates", mock.Anything, req).Return([]string{"first", "second", "third"}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				IsUnitTest:        true,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet("data.akamai_iam_states.test", "id"),
							resource.TestCheckResourceAttr("data.akamai_iam_states.test", "states.#", "3"),
							resource.TestCheckResourceAttr("data.akamai_iam_states.test", "states.0", "first"),
							resource.TestCheckResourceAttr("data.akamai_iam_states.test", "states.1", "second"),
							resource.TestCheckResourceAttr("data.akamai_iam_states.test", "states.2", "third"),
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

		req := iam.ListStatesRequest{Country: "test country"}
		client.On("ListStates", mock.Anything, req).Return([]string{}, errors.New("Could not get states"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				IsUnitTest:        true,
				Steps: []resource.TestStep{
					{
						Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
						ExpectError: regexp.MustCompile(`Could not get states`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
