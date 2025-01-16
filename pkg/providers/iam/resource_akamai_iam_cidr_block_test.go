package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
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
	testCIDR = commonDataForResource{
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
	}

	testCIDRNoComments = commonDataForResource{
		cidrBlock: "128.5.6.5/24",
		enabled:   false,
		actions: &action{
			deleteAction: true,
			editAction:   true,
		},
		cidrBlockID:  1111,
		createdBy:    "jdoe",
		createdDate:  "2006-01-02T15:04:05.999999999Z",
		modifiedBy:   "jkowalski",
		modifiedDate: "2006-01-02T15:04:05.999999999Z",
	}

	updatedCIDR = commonDataForResource{
		cidrBlock: "128.5.6.99/24",
		enabled:   false,
		comments:  ptr.To("test-updated"),
		actions: &action{
			deleteAction: true,
			editAction:   true,
		},
		cidrBlockID:  1111,
		createdBy:    "jdoe",
		createdDate:  "2006-01-02T15:04:05.999999999Z",
		modifiedBy:   "jkowalski",
		modifiedDate: "2006-01-02T15:04:05.999999999Z",
	}
)

func TestCIDRBlockResource(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		configPath string
		init       func(*iam.Mock, commonDataForResource, commonDataForResource)
		createData commonDataForResource
		updateData commonDataForResource
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"happy path - create with comment": {
			init: func(m *iam.Mock, createData, _ commonDataForResource) {
				// Create
				mockCreateCIDRBlock(t, m, createData)
				// Read
				mockGetCIDRBlock(t, m, createData).Twice()
				// Delete
				mockDeleteCIDRBlock(m, createData)
			},
			createData: testCIDR,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.delete", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.edit", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "comments", "test"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block", "128.5.6.5/24"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "enabled", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block_id", "1111"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_by", "jdoe"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_date", "2006-01-02T15:04:05.999999999Z"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_by", "jkowalski"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_date", "2006-01-02T15:04:05.999999999Z"),
					),
				},
			},
		},
		"happy path - create without comment": {
			init: func(m *iam.Mock, createData, _ commonDataForResource) {
				// Create
				mockCreateCIDRBlock(t, m, createData)
				// Read
				mockGetCIDRBlock(t, m, createData).Times(2)
				// Delete
				mockDeleteCIDRBlock(m, createData)
			},
			createData: testCIDRNoComments,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/create_without_comments.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.delete", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.edit", "true"),
						resource.TestCheckNoResourceAttr("akamai_iam_cidr_block.test", "comments"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block", "128.5.6.5/24"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "enabled", "false"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block_id", "1111"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_by", "jdoe"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_date", "2006-01-02T15:04:05.999999999Z"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_by", "jkowalski"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_date", "2006-01-02T15:04:05.999999999Z"),
					),
				},
			},
		},
		"happy path - update all fields": {
			init: func(m *iam.Mock, createData, updateData commonDataForResource) {
				// Create
				mockCreateCIDRBlock(t, m, createData)
				mockGetCIDRBlock(t, m, createData)
				// Read
				mockGetCIDRBlock(t, m, createData).Twice()
				// Update
				mockUpdateCIDRBlock(t, m, updateData)
				// Read
				mockGetCIDRBlock(t, m, updateData)
				// Delete
				mockDeleteCIDRBlock(m, createData)
			},
			createData: testCIDR,
			updateData: updatedCIDR,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.delete", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.edit", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "comments", "test"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block", "128.5.6.5/24"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "enabled", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block_id", "1111"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_by", "jdoe"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_date", "2006-01-02T15:04:05.999999999Z"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_by", "jkowalski"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_date", "2006-01-02T15:04:05.999999999Z"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.delete", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.edit", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "comments", "test-updated"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block", "128.5.6.99/24"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "enabled", "false"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block_id", "1111"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_by", "jdoe"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_date", "2006-01-02T15:04:05.999999999Z"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_by", "jkowalski"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_date", "2006-01-02T15:04:05.999999999Z"),
					),
				},
			},
		},
		"validation error - missing cidr block": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/missing_cidr_block.tf"),
					ExpectError: regexp.MustCompile(`The argument "cidr_block" is required, but no definition was found`),
				},
			},
		},
		"validation error - missing enabled": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/missing_enabled.tf"),
					ExpectError: regexp.MustCompile(`The argument "enabled" is required, but no definition was found`),
				},
			},
		},
		"expect error - create": {
			init: func(m *iam.Mock, createData, _ commonDataForResource) {
				m.On("CreateCIDRBlock", testutils.MockContext, iam.CreateCIDRBlockRequest{
					CIDRBlock: createData.cidrBlock,
					Comments:  createData.comments,
					Enabled:   createData.enabled,
				}).Return(nil, fmt.Errorf("create failed")).Once()
			},
			createData: testCIDR,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/create.tf"),
					ExpectError: regexp.MustCompile(`create failed`),
				},
			},
		},
		"expect error - delete": {
			init: func(m *iam.Mock, createData, _ commonDataForResource) {
				// Create
				mockCreateCIDRBlock(t, m, createData)
				// Read
				mockGetCIDRBlock(t, m, createData).Twice()
				mockGetCIDRBlock(t, m, createData)
				// Delete - error
				m.On("DeleteCIDRBlock", testutils.MockContext, iam.DeleteCIDRBlockRequest{
					CIDRBlockID: createData.cidrBlockID,
				}).Return(iam.ErrDeleteCIDRBlock).Once()
				// Delete - destroy
				mockDeleteCIDRBlock(m, createData)
			},
			createData: testCIDR,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.delete", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "actions.edit", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "comments", "test"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block", "128.5.6.5/24"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "enabled", "true"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "cidr_block_id", "1111"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_by", "jdoe"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "created_date", "2006-01-02T15:04:05.999999999Z"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_by", "jkowalski"),
						resource.TestCheckResourceAttr("akamai_iam_cidr_block.test", "modified_date", "2006-01-02T15:04:05.999999999Z"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/empty_config.tf"),
					ExpectError: regexp.MustCompile(`Error: delete cidr block {2 1111} failed`),
				},
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client, tc.createData, tc.updateData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}

}

func TestImportCIDRBlockResource(t *testing.T) {
	tests := map[string]struct {
		importID    string
		init        func(*iam.Mock, commonDataForResource)
		mockData    commonDataForResource
		expectError *regexp.Regexp
	}{
		"happy path - import with comments": {
			importID: "1111",
			mockData: testCIDR,
			init: func(m *iam.Mock, data commonDataForResource) {
				// Read
				mockGetCIDRBlock(t, m, data)
			},
		},
		"happy path - import without comments": {
			importID: "1111",
			mockData: testCIDRNoComments,
			init: func(m *iam.Mock, data commonDataForResource) {
				// Read
				mockGetCIDRBlock(t, m, data)
			},
		},
		"expect error - wrong import ID": {
			importID:    "wrong format",
			expectError: regexp.MustCompile(`Error: could not convert import ID to int`),
		},
		"expect error - read": {
			importID: "1111",
			init: func(m *iam.Mock, data commonDataForResource) {
				m.On("GetCIDRBlock", testutils.MockContext, iam.GetCIDRBlockRequest{
					CIDRBlockID: data.cidrBlockID,
					Actions:     true,
				}).Return(nil, fmt.Errorf("get failed")).Once()
			},
			mockData:    testCIDR,
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: checkImportCIDRBlock(tc.mockData),
							ImportStateId:    tc.importID,
							ImportState:      true,
							ResourceName:     "akamai_iam_cidr_block.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResCIDRBlock/importable.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkImportCIDRBlock(data commonDataForResource) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		attr := make(map[string]string)

		if len(s) == 0 {
			return errors.New("no instance found")
		}
		if len(s) != 1 {
			return fmt.Errorf("expected one instance: %d", len(s))
		}

		if data.actions != nil {
			attr["actions.delete"] = strconv.FormatBool(data.actions.deleteAction)
			attr["actions.edit"] = strconv.FormatBool(data.actions.editAction)
		}
		if data.comments != nil {
			attr["comments"] = *data.comments
		} else {
			attr["comments"] = ""
		}
		attr["cidr_block"] = data.cidrBlock
		attr["enabled"] = strconv.FormatBool(data.enabled)
		attr["cidr_block_id"] = strconv.FormatInt(data.cidrBlockID, 10)
		attr["created_by"] = data.createdBy
		attr["created_date"] = data.createdDate
		attr["modified_by"] = data.modifiedBy
		attr["modified_date"] = data.modifiedDate

		state := s[0].Attributes

		attributes := attr

		var invalidValues []string
		for field, expectedVal := range attributes {
			if state[field] != expectedVal {
				invalidValues = append(invalidValues, fmt.Sprintf("field: %s, got: %s, expected: %s ", field, state[field], expectedVal))
			}
		}

		if len(invalidValues) > 0 {
			return fmt.Errorf("found invalid values: %s", strings.Join(invalidValues, "\n"))
		}
		return nil
	}
}

func mockCreateCIDRBlock(t *testing.T, m *iam.Mock, testData commonDataForResource) *mock.Call {
	var act *iam.CIDRActions

	if testData.actions != nil {
		act = &iam.CIDRActions{
			Delete: testData.actions.deleteAction,
			Edit:   testData.actions.editAction,
		}
	}

	return m.On("CreateCIDRBlock", testutils.MockContext, iam.CreateCIDRBlockRequest{
		CIDRBlock: testData.cidrBlock,
		Comments:  testData.comments,
		Enabled:   testData.enabled,
	}).Return(&iam.CreateCIDRBlockResponse{
		Actions:      act,
		CIDRBlock:    testData.cidrBlock,
		CIDRBlockID:  testData.cidrBlockID,
		Comments:     testData.comments,
		CreatedBy:    testData.createdBy,
		CreatedDate:  test.NewTimeFromString(t, testData.createdDate),
		Enabled:      testData.enabled,
		ModifiedBy:   testData.modifiedBy,
		ModifiedDate: test.NewTimeFromString(t, testData.modifiedDate),
	}, nil).Once()
}

func mockGetCIDRBlock(t *testing.T, m *iam.Mock, testData commonDataForResource) *mock.Call {
	var act *iam.CIDRActions

	if testData.actions != nil {
		act = &iam.CIDRActions{
			Delete: testData.actions.deleteAction,
			Edit:   testData.actions.editAction,
		}
	}

	return m.On("GetCIDRBlock", testutils.MockContext, iam.GetCIDRBlockRequest{
		CIDRBlockID: testData.cidrBlockID,
		Actions:     true,
	}).Return(&iam.GetCIDRBlockResponse{
		Actions:      act,
		CIDRBlock:    testData.cidrBlock,
		CIDRBlockID:  testData.cidrBlockID,
		Comments:     testData.comments,
		CreatedBy:    testData.createdBy,
		CreatedDate:  test.NewTimeFromString(t, testData.createdDate),
		Enabled:      testData.enabled,
		ModifiedBy:   testData.modifiedBy,
		ModifiedDate: test.NewTimeFromString(t, testData.modifiedDate),
	}, nil).Once()
}

func mockDeleteCIDRBlock(m *iam.Mock, testData commonDataForResource) *mock.Call {
	return m.On("DeleteCIDRBlock", testutils.MockContext, iam.DeleteCIDRBlockRequest{
		CIDRBlockID: testData.cidrBlockID,
	}).Return(nil).Once()
}

func mockUpdateCIDRBlock(t *testing.T, m *iam.Mock, testData commonDataForResource) *mock.Call {
	var act *iam.CIDRActions

	if testData.actions != nil {
		act = &iam.CIDRActions{
			Delete: testData.actions.deleteAction,
			Edit:   testData.actions.editAction,
		}
	}

	return m.On("UpdateCIDRBlock", testutils.MockContext, iam.UpdateCIDRBlockRequest{
		CIDRBlockID: testData.cidrBlockID,
		Body: iam.UpdateCIDRBlockRequestBody{
			CIDRBlock: testData.cidrBlock,
			Comments:  testData.comments,
			Enabled:   testData.enabled,
		},
	}).Return(&iam.UpdateCIDRBlockResponse{
		Actions:      act,
		CIDRBlock:    testData.cidrBlock,
		CIDRBlockID:  testData.cidrBlockID,
		Comments:     testData.comments,
		CreatedBy:    testData.createdBy,
		CreatedDate:  test.NewTimeFromString(t, testData.createdDate),
		Enabled:      testData.enabled,
		ModifiedBy:   testData.modifiedBy,
		ModifiedDate: test.NewTimeFromString(t, testData.modifiedDate),
	}, nil).Once()
}
