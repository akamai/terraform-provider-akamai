package iam

import (
	"testing"

	iam "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDSStates(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		client := &IAM{}
		client.Test(test.TattleT{T: t})

		req := iam.ListStatesRequest{Country: "test country"}
		client.On("ListStates", mock.Anything, req).Return([]string{"first", "second", "third"}, nil)

		p := provider{}
		p.SetCache(metaCache{})
		p.SetClient(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s/step0.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_states.test", "id"),
						resource.TestCheckTypeSetElemAttr("data.akamai_iam_states.test", "states.*", "first"),
						resource.TestCheckTypeSetElemAttr("data.akamai_iam_states.test", "states.*", "second"),
						resource.TestCheckTypeSetElemAttr("data.akamai_iam_states.test", "states.*", "third"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
