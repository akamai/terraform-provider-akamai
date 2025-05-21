package edgeworkers

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataEdgeworkersResourceTier(t *testing.T) {
	tests := map[string]struct {
		init  func(*edgeworkers.Mock)
		steps []resource.TestStep
	}{
		"read resource tier": {
			init: func(m *edgeworkers.Mock) {
				m.On("ListResourceTiers", testutils.MockContext, edgeworkers.ListResourceTiersRequest{
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/basic.tf"),
					Check:  resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "resource_tier_id", "100"),
				},
			},
		},
		"ctr contract prefix": {
			init: func(m *edgeworkers.Mock) {
				m.On("ListResourceTiers", testutils.MockContext, edgeworkers.ListResourceTiersRequest{
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/ctr_prefix.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "resource_tier_id", "100"),
						resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "contract_id", "ctr_1-599K"),
					),
				},
			},
		},
		"ctr contract prefix update": {
			init: func(m *edgeworkers.Mock) {
				m.On("ListResourceTiers", testutils.MockContext, edgeworkers.ListResourceTiersRequest{
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "resource_tier_id", "100"),
						resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "contract_id", "1-599K"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/ctr_prefix.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "resource_tier_id", "100"),
						resource.TestCheckResourceAttr("data.akamai_edgeworkers_resource_tier.test", "contract_id", "ctr_1-599K"),
					),
				},
			},
		},
		"contract id not exist": {
			init: func(m *edgeworkers.Mock) {
				m.On("ListResourceTiers", testutils.MockContext, edgeworkers.ListResourceTiersRequest{
					ContractID: "1-599K",
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/basic.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"resource tier name not exist": {
			init: func(m *edgeworkers.Mock) {
				m.On("ListResourceTiers", testutils.MockContext, edgeworkers.ListResourceTiersRequest{
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
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/incorrect_resource_tier_name.tf"),
					ExpectError: regexp.MustCompile("Resource tier with name: 'Incorrect' was not found"),
				},
			},
		},
		"missing constract id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/missing_contract_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
			},
		},
		"missing resource tier name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDataEdgeWorkersResourceTier/missing_resource_tier_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "resource_tier_name" is required, but no definition was found`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
