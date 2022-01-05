package edgeworkers

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataEdgeworkersResourceTier(t *testing.T) {
	tests := map[string]struct {
		init  func(*mockedgeworkers)
		steps []resource.TestStep
	}{
		"read resource tier": {
			init: func(m *mockedgeworkers) {
				m.On("ListResourceTiers", mock.Anything, edgeworkers.ListResourceTiersRequest{
					ContractID: "1-599K",
				}).Return(&edgeworkers.ListResourceTiersResponse{
					ResourceTiers: []edgeworkers.ResourceTier{
						{
							ID:   100,
							Name: "Basic Compute",
						},
						{
							ID:   200,
							Name: "Dynamic Compute",
						},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestDataEdgeworkersResourceTier/basic.tf"),
					Check:  resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "resource_tier_id", "100"),
				},
			},
		},
		"contract id not exist": {
			init: func(m *mockedgeworkers) {
				m.On("ListResourceTiers", mock.Anything, edgeworkers.ListResourceTiersRequest{
					ContractID: "1-599K",
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestDataEdgeworkersResourceTier/basic.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"resource tier name not exist": {
			init: func(m *mockedgeworkers) {
				m.On("ListResourceTiers", mock.Anything, edgeworkers.ListResourceTiersRequest{
					ContractID: "1-599K",
				}).Return(&edgeworkers.ListResourceTiersResponse{
					ResourceTiers: []edgeworkers.ResourceTier{
						{
							ID:   100,
							Name: "Basic Compute",
						},
						{
							ID:   200,
							Name: "Dynamic Compute",
						},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestDataEdgeworkersResourceTier/incorrect_resource_tier_name.tf"),
					ExpectError: regexp.MustCompile("Resource tier with name: 'Incorrect' was not found"),
				},
			},
		},
		"missing constract id": {
			init: func(m *mockedgeworkers) {},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestDataEdgeworkersResourceTier/missing_contract_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
			},
		},
		"missing resource tier name": {
			init: func(m *mockedgeworkers) {},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestDataEdgeworkersResourceTier/missing_resource_tier_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "resource_tier_name" is required, but no definition was found`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockedgeworkers{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
