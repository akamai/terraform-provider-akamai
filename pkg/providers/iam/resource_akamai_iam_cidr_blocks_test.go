package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/internal/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"
)

type (
	commonDataForResource struct {
		cidrBlock    string
		comments     *string
		enabled      bool
		actions      *action
		cidrBlockID  int64
		createdBy    string
		createdDate  string
		modifiedBy   string
		modifiedDate string
	}

	action struct {
		deleteAction bool
		editAction   bool
	}
)

var (
	createdCidrBlocks = []commonDataForResource{
		{
			cidrBlock: "128.5.6.5/24",
			enabled:   true,
			comments:  ptr.To("test"),
			actions: &action{
				deleteAction: true,
				editAction:   true,
			},
			cidrBlockID:  1111,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.5.6.6/24",
			enabled:   false,
			comments:  nil,
			actions: &action{
				deleteAction: true,
				editAction:   false,
			},
			cidrBlockID:  2222,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.5.6.7/24",
			comments:  ptr.To("test1234"),
			enabled:   true,
			actions: &action{
				deleteAction: false,
				editAction:   true,
			},
			cidrBlockID:  3333,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.5.6.8/24",
			comments:  ptr.To("abcd12345"),
			enabled:   false,
			actions: &action{
				deleteAction: true,
				editAction:   true,
			},
			cidrBlockID:  4444,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.5.6.9/24",
			enabled:   true,
			comments:  nil,
			actions: &action{
				deleteAction: false,
				editAction:   false,
			},
			cidrBlockID:  5555,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
	}

	updatedCidrBlocks = []commonDataForResource{
		{
			cidrBlock: "128.1.2.5/24",
			enabled:   false,
			comments:  nil,
			actions: &action{
				deleteAction: true,
				editAction:   true,
			},
			cidrBlockID:  1111,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.1.2.6/24",
			enabled:   false,
			comments:  nil,
			actions: &action{
				deleteAction: true,
				editAction:   false,
			},
			cidrBlockID:  2222,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.1.2.7/24",
			comments:  ptr.To("test1234"),
			enabled:   true,
			actions: &action{
				deleteAction: false,
				editAction:   true,
			},
			cidrBlockID:  3333,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.1.2.8/24",
			comments:  ptr.To("up12345"),
			enabled:   true,
			actions: &action{
				deleteAction: true,
				editAction:   true,
			},
			cidrBlockID:  4444,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.1.2.9/24",
			enabled:   false,
			comments:  nil,
			actions: &action{
				deleteAction: false,
				editAction:   false,
			},
			cidrBlockID:  5555,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
	}

	threeOutOffiveUpdatedCidrBlocks = []commonDataForResource{
		{
			cidrBlock: "128.2.2.5/28",
			enabled:   false,
			comments:  nil,
			actions: &action{
				deleteAction: true,
				editAction:   true,
			},
			cidrBlockID:  1111,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.2.2.6/28",
			enabled:   true,
			comments:  nil,
			actions: &action{
				deleteAction: true,
				editAction:   false,
			},
			cidrBlockID:  2222,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.2.2.7/28",
			comments:  ptr.To("test12345"),
			enabled:   false,
			actions: &action{
				deleteAction: false,
				editAction:   true,
			},
			cidrBlockID:  3333,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.5.6.8/24",
			comments:  ptr.To("abcd12345"),
			enabled:   false,
			actions: &action{
				deleteAction: true,
				editAction:   true,
			},
			cidrBlockID:  4444,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
		{
			cidrBlock: "128.5.6.9/24",
			enabled:   true,
			comments:  nil,
			actions: &action{
				deleteAction: false,
				editAction:   false,
			},
			cidrBlockID:  5555,
			createdBy:    "jdoe",
			createdDate:  "2006-01-02T15:04:05.999999999Z",
			modifiedBy:   "jkowalski",
			modifiedDate: "2006-01-02T15:04:05.999999999Z",
		},
	}
)

func TestResourceIAMCIDRBlocksResource(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		configPath  string
		init        func(*testing.T, *iam.Mock, []commonDataForResource, []commonDataForResource)
		createData  []commonDataForResource
		updatedData []commonDataForResource
		steps       []resource.TestStep
		error       *regexp.Regexp
	}{
		"create - single cidr Block": {
			init: func(t *testing.T, m *iam.Mock, resourceData []commonDataForResource, _ []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, resourceData)

				// Read
				mockListCIDRBlock(t, m, resourceData)

				// Delete
				mockDeleteCIDRBlocks(m, resourceData)
			},
			createData: generateCIDRBlocks(createdCidrBlocks[:1]),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create.tf"),
					Check:  checkAttrs(generateCIDRBlocks(createdCidrBlocks[:1])),
				},
			},
		},
		"create - multiple cidr Blocks": {
			init: func(t *testing.T, m *iam.Mock, resourceData []commonDataForResource, _ []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, resourceData)

				// Read
				mockListCIDRBlock(t, m, resourceData)

				// Delete
				mockDeleteCIDRBlocks(m, resourceData)
			},
			createData: generateCIDRBlocks(createdCidrBlocks),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create_multiple.tf"),
					Check:  checkAttrs(generateCIDRBlocks(createdCidrBlocks)),
				},
			},
		},
		"update - single cidr block": {
			init: func(t *testing.T, m *iam.Mock, createData, updateData []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, createData)

				// Read
				mockListCIDRBlock(t, m, createData)

				// Refresh read
				mockListCIDRBlock(t, m, createData)

				// Update
				mockUpdateCIDRBlocks(t, m, updateData)

				// Read
				mockListCIDRBlock(t, m, updateData)

				// Delete
				mockDeleteCIDRBlocks(m, updateData)
			},
			createData:  generateCIDRBlocks(createdCidrBlocks[:1]),
			updatedData: generateCIDRBlocks(updatedCidrBlocks[:1]),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create.tf"),
					Check:  checkAttrs(generateCIDRBlocks(createdCidrBlocks[:1])),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/update_single.tf"),
					Check:  checkAttrs(generateCIDRBlocks(updatedCidrBlocks[:1])),
				},
			},
		},
		"update - all cidr blocks": {
			init: func(t *testing.T, m *iam.Mock, createData, updateData []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, createData)

				// Read
				mockListCIDRBlock(t, m, createData)

				// Refresh read
				mockListCIDRBlock(t, m, createData)

				// Update
				mockUpdateCIDRBlocks(t, m, updateData)

				// Read
				mockListCIDRBlock(t, m, updateData)

				// Delete
				mockDeleteCIDRBlocks(m, updateData)
			},
			createData:  generateCIDRBlocks(createdCidrBlocks),
			updatedData: generateCIDRBlocks(updatedCidrBlocks),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create_multiple.tf"),
					Check:  checkAttrs(createdCidrBlocks),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/update_five.tf"),
					Check:  checkAttrs(updatedCidrBlocks),
				},
			},
		},
		"update - 3 cidr blocks from 5 created": {
			init: func(t *testing.T, m *iam.Mock, createData, updateData []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, createData)

				// Read
				mockListCIDRBlock(t, m, createData)

				// Refresh read
				mockListCIDRBlock(t, m, createData)

				// Update
				mockUpdateCIDRBlocks(t, m, updateData)

				// Read
				mockListCIDRBlock(t, m, updateData)

				// Delete
				mockDeleteCIDRBlocks(m, updateData)
			},
			createData:  generateCIDRBlocks(createdCidrBlocks),
			updatedData: generateCIDRBlocks(threeOutOffiveUpdatedCidrBlocks),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create_multiple.tf"),
					Check:  checkAttrs(createdCidrBlocks),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/update_three_out_of_five.tf"),
					Check:  checkAttrs(threeOutOffiveUpdatedCidrBlocks),
				},
			},
		},
		"update - 3 cidr blocks and 2 remove": {
			init: func(t *testing.T, m *iam.Mock, createData, updateData []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, createData)

				// Read
				mockListCIDRBlock(t, m, createData)

				// Refresh read
				mockListCIDRBlock(t, m, createData)

				// Update
				mockUpdateCIDRBlocks(t, m, updateData)

				// Read
				mockListCIDRBlock(t, m, updateData)

				// Delete
				mockDeleteCIDRBlocks(m, updateData)
			},
			createData:  generateCIDRBlocks(createdCidrBlocks),
			updatedData: generateCIDRBlocks(updatedCidrBlocks[:3]),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create_multiple.tf"),
					Check:  checkAttrs(createdCidrBlocks),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/update_three_remove_two.tf"),
					Check:  checkAttrs(generateCIDRBlocks(updatedCidrBlocks[:3])),
				},
			},
		},
		"missing cidr block": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/missing_cidr_block.tf"),
				ExpectError: regexp.MustCompile("\\s*Inappropriate value for attribute \"cidr_blocks\": element 0: attribute\\s*\"cidr_block\" is required."),
			}},
		},
		"missing enabled": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/missing_enabled.tf"),
				ExpectError: regexp.MustCompile("\\s*Inappropriate value for attribute \"cidr_blocks\": element 0: attribute\\s*\"enabled\" is required."),
			}},
		},
		"error - create": {
			init: func(t *testing.T, m *iam.Mock, resourceData, _ []commonDataForResource) {
				m.On("CreateCIDRBlock", mock.Anything, iam.CreateCIDRBlockRequest{
					CIDRBlock: resourceData[0].cidrBlock,
					Comments:  resourceData[0].comments,
					Enabled:   resourceData[0].enabled,
				}).Return(nil, iam.ErrCreateCIDRBlock).Once()
			},
			createData: generateCIDRBlocks(createdCidrBlocks[:1]),
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create.tf"),
				ExpectError: regexp.MustCompile("create cidr block failed"),
			}},
		},
		"error - delete": {
			init: func(t *testing.T, m *iam.Mock, resourceData, _ []commonDataForResource) {
				// Create
				mockCreateCIDRBlocks(t, m, resourceData)

				// Read
				mockListCIDRBlock(t, m, resourceData)

				// Read
				mockListCIDRBlock(t, m, resourceData)

				// Delete - error
				m.On("DeleteCIDRBlock", mock.Anything, iam.DeleteCIDRBlockRequest{
					CIDRBlockID: resourceData[0].cidrBlockID,
				}).Return(iam.ErrDeleteCIDRBlock).Once()

				// Delete - destroy
				mockDeleteCIDRBlocks(m, resourceData)
			},
			createData: generateCIDRBlocks(createdCidrBlocks[:1]),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/create.tf"),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/empty_config.tf"),
					ExpectError: regexp.MustCompile(`Error: delete cidr block {2 1111} failed`),
				},
			},
		},
	}
	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &iam.Mock{}
			if test.init != nil {
				test.init(t, client, test.createData, test.updatedData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}

}

func TestImportIAMCIDRBlocksResource(t *testing.T) {
	tests := map[string]struct {
		importID   string
		configPath string
		init       func(*testing.T, *iam.Mock, []commonDataForResource)
		mockData   []commonDataForResource
	}{
		"import - single cidr block": {
			importID: " ",
			mockData: generateCIDRBlocks(createdCidrBlocks[:1]),
			init: func(t *testing.T, m *iam.Mock, resourceData []commonDataForResource) {
				// Import
				mockListCIDRBlock(t, m, resourceData)

				// Read
				mockListCIDRBlock(t, m, resourceData)
			},
		},
		"import - five cidr blocks": {
			importID: " ",
			mockData: generateCIDRBlocks(createdCidrBlocks),
			init: func(t *testing.T, m *iam.Mock, resourceData []commonDataForResource) {
				// Import
				mockListCIDRBlock(t, m, resourceData)

				// Read
				mockListCIDRBlock(t, m, resourceData)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			test.init(t, client, test.mockData)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: checkImportCIDRBlocks(test.mockData),
							ImportStateId:    test.importID,
							ImportState:      true,
							ResourceName:     "akamai_iam_cidr_blocks.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResCIDRBlocks/importable.tf"),
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkImportCIDRBlocks(data []commonDataForResource) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		attr := make(map[string]string)

		if len(s) == 0 {
			return errors.New("No Instance found")
		}
		if len(s) != 1 {
			return fmt.Errorf("Expected one Instance: %d", len(s))
		}

		for i, cidr := range data {
			if cidr.actions != nil {
				attr[fmt.Sprintf("cidr_blocks.%d.actions.delete", i)] = strconv.FormatBool(cidr.actions.deleteAction)
				attr[fmt.Sprintf("cidr_blocks.%d.actions.edit", i)] = strconv.FormatBool(cidr.actions.editAction)
			}
			if cidr.comments != nil {
				attr[fmt.Sprintf("cidr_blocks.%d.comments", i)] = *cidr.comments
			} else {
				attr[fmt.Sprintf("cidr_blocks.%d.comments", i)] = ""
			}
			attr[fmt.Sprintf("cidr_blocks.%d.cidr_block", i)] = cidr.cidrBlock
			attr[fmt.Sprintf("cidr_blocks.%d.enabled", i)] = strconv.FormatBool(cidr.enabled)
			attr[fmt.Sprintf("cidr_blocks.%d.cidr_block_id", i)] = strconv.FormatInt(cidr.cidrBlockID, 10)
			attr[fmt.Sprintf("cidr_blocks.%d.created_by", i)] = cidr.createdBy
			attr[fmt.Sprintf("cidr_blocks.%d.created_date", i)] = cidr.createdDate
			attr[fmt.Sprintf("cidr_blocks.%d.modified_by", i)] = cidr.modifiedBy
			attr[fmt.Sprintf("cidr_blocks.%d.modified_date", i)] = cidr.modifiedDate
		}

		state := s[0].Attributes

		attributes := attr

		invalidValues := []string{}
		for field, expectedVal := range attributes {
			if state[field] != expectedVal {
				invalidValues = append(invalidValues, fmt.Sprintf("field: %s, got: %s, expected: %s ", field, state[field], expectedVal))
			}
		}

		if len(invalidValues) != 0 {
			return fmt.Errorf(strings.Join(invalidValues, "\n"))
		}
		return nil
	}
}

func mockCreateCIDRBlocks(t *testing.T, m *iam.Mock, testData []commonDataForResource) []*mock.Call {
	var mocks []*mock.Call
	var act *iam.CIDRActions

	for _, d := range testData {
		if d.actions != nil {
			act = &iam.CIDRActions{
				Delete: d.actions.deleteAction,
				Edit:   d.actions.editAction,
			}
		}
		mocks = append(mocks, m.On("CreateCIDRBlock", mock.Anything, iam.CreateCIDRBlockRequest{
			CIDRBlock: d.cidrBlock,
			Comments:  d.comments,
			Enabled:   d.enabled,
		}).Return(&iam.CreateCIDRBlockResponse{
			Actions:      act,
			CIDRBlock:    d.cidrBlock,
			CIDRBlockID:  d.cidrBlockID,
			Comments:     d.comments,
			CreatedBy:    d.createdBy,
			CreatedDate:  test.NewTimeFromString(t, d.createdDate),
			Enabled:      d.enabled,
			ModifiedBy:   d.modifiedBy,
			ModifiedDate: test.NewTimeFromString(t, d.modifiedDate),
		}, nil).Once())
	}

	return mocks
}

func mockListCIDRBlock(t *testing.T, m *iam.Mock, testData []commonDataForResource) *mock.Call {
	var cidrBlock iam.ListCIDRBlocksResponse
	var act *iam.CIDRActions

	for _, cidr := range testData {
		if cidr.actions != nil {
			act = &iam.CIDRActions{
				Delete: cidr.actions.deleteAction,
				Edit:   cidr.actions.editAction,
			}
		}
		cidrBlock = append(cidrBlock, iam.CIDRBlock{
			Actions:      act,
			CIDRBlock:    cidr.cidrBlock,
			CIDRBlockID:  cidr.cidrBlockID,
			Comments:     cidr.comments,
			CreatedBy:    cidr.createdBy,
			CreatedDate:  test.NewTimeFromString(t, cidr.createdDate),
			Enabled:      cidr.enabled,
			ModifiedBy:   cidr.modifiedBy,
			ModifiedDate: test.NewTimeFromString(t, cidr.modifiedDate),
		})
	}

	return m.On("ListCIDRBlocks", mock.Anything, iam.ListCIDRBlocksRequest{
		Actions: true,
	}).Return(cidrBlock, nil).Once()
}

func mockDeleteCIDRBlocks(m *iam.Mock, testData []commonDataForResource) []*mock.Call {
	var mocks []*mock.Call

	for _, d := range testData {
		mocks = append(mocks, m.On("DeleteCIDRBlock", mock.Anything, iam.DeleteCIDRBlockRequest{
			CIDRBlockID: d.cidrBlockID,
		}).Return(nil).Once())
	}

	return mocks
}

func mockUpdateCIDRBlocks(t *testing.T, m *iam.Mock, data []commonDataForResource) []*mock.Call {
	var mocks []*mock.Call

	for _, d := range data {
		mocks = append(mocks, m.On("UpdateCIDRBlock", mock.Anything, iam.UpdateCIDRBlockRequest{
			CIDRBlockID: d.cidrBlockID,
			Body: iam.UpdateCIDRBlockRequestBody{
				CIDRBlock: d.cidrBlock,
				Comments:  d.comments,
				Enabled:   d.enabled,
			},
		}).Return(&iam.UpdateCIDRBlockResponse{
			Actions: &iam.CIDRActions{
				Delete: d.actions.deleteAction,
				Edit:   d.actions.editAction,
			},
			CIDRBlock:    d.cidrBlock,
			CIDRBlockID:  d.cidrBlockID,
			Comments:     d.comments,
			CreatedBy:    d.createdBy,
			CreatedDate:  test.NewTimeFromString(t, d.createdDate),
			Enabled:      d.enabled,
			ModifiedBy:   d.modifiedBy,
			ModifiedDate: test.NewTimeFromString(t, d.modifiedDate),
		}, nil).Once())
	}

	return mocks
}

func checkAttrs(data []commonDataForResource) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	for i, cidr := range data {
		if cidr.actions != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.actions.delete", i), strconv.FormatBool(cidr.actions.deleteAction)))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.actions.edit", i), strconv.FormatBool(cidr.actions.editAction)))
		}
		if cidr.comments != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttrPtr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.comments", i), cidr.comments))
		} else {
			checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.comments", i)))
		}
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.cidr_block", i), cidr.cidrBlock))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.enabled", i), strconv.FormatBool(cidr.enabled)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.cidr_block_id", i), strconv.FormatInt(cidr.cidrBlockID, 10)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.created_by", i), cidr.createdBy))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.created_date", i), cidr.createdDate))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.modified_by", i), cidr.modifiedBy))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_iam_cidr_blocks.test", fmt.Sprintf("cidr_blocks.%d.modified_date", i), cidr.modifiedDate))
	}
	return resource.ComposeAggregateTestCheckFunc(
		checkFuncs...,
	)
}

func generateCIDRBlocks(data []commonDataForResource) []commonDataForResource {
	var cidrBlocks []commonDataForResource

	for _, cidr := range data {
		cidrBlocks = append(cidrBlocks, commonDataForResource{
			cidrBlock:    cidr.cidrBlock,
			comments:     cidr.comments,
			enabled:      cidr.enabled,
			actions:      cidr.actions,
			cidrBlockID:  cidr.cidrBlockID,
			createdBy:    cidr.createdBy,
			createdDate:  cidr.createdDate,
			modifiedBy:   cidr.modifiedBy,
			modifiedDate: cidr.modifiedDate,
		})
	}

	return cidrBlocks
}
