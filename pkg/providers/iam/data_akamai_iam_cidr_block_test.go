package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/internal/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	testDataForCIDRBlock struct {
		cidrBlockID  int64
		actions      *testActions
		cidrBlock    string
		comments     string
		createdBy    string
		createdDate  string
		enabled      bool
		modifiedBy   string
		modifiedDate string
	}

	testActions struct {
		delete bool
		edit   bool
	}
)

var (
	basicTestDataForCIDRBlock = testDataForCIDRBlock{
		cidrBlockID: 2567,
		actions: &testActions{
			delete: true,
			edit:   false,
		},
		cidrBlock:    "128.5.6.6/24",
		comments:     "APAC Region",
		createdBy:    "alfulani",
		createdDate:  "2017-07-27T18:11:25Z",
		enabled:      true,
		modifiedBy:   "alfulani",
		modifiedDate: "2017-07-27T18:11:25Z",
	}
)

func TestDataCIDRBlock(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *iam.Mock, testDataForCIDRBlock)
		mockData   testDataForCIDRBlock
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataCIDRBlock/default.tf",
			init: func(t *testing.T, m *iam.Mock, mockData testDataForCIDRBlock) {
				expectGetCIDRBlock(t, m, mockData, 3)
			},
			mockData: basicTestDataForCIDRBlock,
		},
		"error - missing cidr_block_id": {
			configPath: "testdata/TestDataCIDRBlock/missing_cidr_block_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
			mockData:   basicTestDataForCIDRBlock,
		},
		"error - GetCIDRBlock call failed ": {
			configPath: "testdata/TestDataCIDRBlock/default.tf",
			init: func(t *testing.T, m *iam.Mock, mockData testDataForCIDRBlock) {
				getCIDRBlockReq := iam.GetCIDRBlockRequest{CIDRBlockID: mockData.cidrBlockID, Actions: true}
				m.On("GetCIDRBlock", mock.Anything, getCIDRBlockReq).Return(nil, errors.New("test error"))
			},
			mockData: basicTestDataForCIDRBlock,
			error:    regexp.MustCompile("test error"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if test.init != nil {
				test.init(t, client, test.mockData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkCIDRBlockAttrs(),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkCIDRBlockAttrs() resource.TestCheckFunc {
	name := "data.akamai_iam_cidr_block.test"

	checksFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "cidr_block_id", "2567"),
		resource.TestCheckResourceAttr(name, "cidr_block", "128.5.6.6/24"),
		resource.TestCheckResourceAttr(name, "comments", "APAC Region"),
		resource.TestCheckResourceAttr(name, "created_by", "alfulani"),
		resource.TestCheckResourceAttr(name, "created_date", "2017-07-27T18:11:25Z"),
		resource.TestCheckResourceAttr(name, "enabled", "true"),
		resource.TestCheckResourceAttr(name, "modified_by", "alfulani"),
		resource.TestCheckResourceAttr(name, "modified_date", "2017-07-27T18:11:25Z"),
	}
	return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
}

func expectGetCIDRBlock(t *testing.T, client *iam.Mock, data testDataForCIDRBlock, timesToRun int) {

	getCIDRBlockReq := iam.GetCIDRBlockRequest{
		CIDRBlockID: data.cidrBlockID,
		Actions:     true,
	}

	createdDate := test.NewTimeFromString(t, data.createdDate)
	modifiedDate := test.NewTimeFromString(t, data.modifiedDate)

	getCIDRBlockResp := iam.GetCIDRBlockResponse{
		CIDRBlock:    data.cidrBlock,
		CIDRBlockID:  data.cidrBlockID,
		Comments:     &data.comments,
		CreatedBy:    data.createdBy,
		CreatedDate:  createdDate,
		Enabled:      data.enabled,
		ModifiedBy:   data.modifiedBy,
		ModifiedDate: modifiedDate,
	}
	if data.actions != nil {
		getCIDRBlockResp.Actions = &iam.CIDRActions{
			Delete: data.actions.delete,
			Edit:   data.actions.edit,
		}
	}
	client.On("GetCIDRBlock", mock.Anything, getCIDRBlockReq).Return(&getCIDRBlockResp, nil).Times(timesToRun)

}
