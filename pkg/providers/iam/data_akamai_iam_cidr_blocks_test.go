package iam

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCIDRBlocks(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*iam.Mock)
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataCIDRBlocks/default.tf",
			init: func(m *iam.Mock) {
				listCIDRBlocksReq := iam.ListCIDRBlocksRequest{Actions: true}
				listCIDRBlocksResp := iam.ListCIDRBlocksResponse{
					{
						Actions: &iam.CIDRActions{
							Delete: true,
							Edit:   false,
						},
						CIDRBlock:    "128.5.6.6/24",
						CIDRBlockID:  2567,
						Comments:     ptr.To("APAC Region"),
						CreatedBy:    "user1",
						CreatedDate:  time.Date(2017, 7, 27, 18, 11, 25, 0, time.UTC),
						Enabled:      true,
						ModifiedBy:   "user1",
						ModifiedDate: time.Date(2017, 7, 27, 18, 11, 25, 0, time.UTC),
					},
					{
						Actions: &iam.CIDRActions{
							Delete: true,
							Edit:   false,
						},
						CIDRBlock:    "128.5.6.6/24",
						CIDRBlockID:  6042,
						Comments:     ptr.To("East Coast Office"),
						CreatedBy:    "user2",
						CreatedDate:  time.Date(2017, 7, 27, 18, 11, 25, 0, time.UTC),
						Enabled:      true,
						ModifiedBy:   "user3",
						ModifiedDate: time.Date(2017, 7, 27, 18, 11, 25, 0, time.UTC),
					},
				}
				m.On("ListCIDRBlocks", testutils.MockContext, listCIDRBlocksReq).Return(listCIDRBlocksResp, nil).Times(3)
			},
		},
		"error - ListCIDRBlocks call failed ": {
			configPath: "testdata/TestDataCIDRBlocks/default.tf",
			init: func(m *iam.Mock) {
				listCIDRBlocksReq := iam.ListCIDRBlocksRequest{Actions: true}
				m.On("ListCIDRBlocks", testutils.MockContext, listCIDRBlocksReq).Return(nil, errors.New("test error"))
			},
			error: regexp.MustCompile("test error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.configPath),
							Check:       checkCIDRBlocksAttrs(),
							ExpectError: tc.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkCIDRBlocksAttrs() resource.TestCheckFunc {
	name := "data.akamai_iam_cidr_blocks.test"
	checksFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "cidr_blocks.#", "2"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.cidr_block_id", "2567"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.cidr_block", "128.5.6.6/24"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.comments", "APAC Region"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.created_by", "user1"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.created_date", "2017-07-27T18:11:25Z"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.enabled", "true"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.modified_by", "user1"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.0.modified_date", "2017-07-27T18:11:25Z"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.cidr_block_id", "6042"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.cidr_block", "128.5.6.6/24"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.comments", "East Coast Office"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.created_by", "user2"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.created_date", "2017-07-27T18:11:25Z"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.enabled", "true"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.modified_by", "user3"),
		resource.TestCheckResourceAttr(name, "cidr_blocks.1.modified_date", "2017-07-27T18:11:25Z"),
	}
	return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
}
