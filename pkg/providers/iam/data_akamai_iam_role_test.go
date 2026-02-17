package iam

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRoleDataSource(t *testing.T) {
	createdDate := time.Date(2017, time.July, 27, 18, 11, 25, 0, time.UTC)
	modifiedDate := time.Date(2017, time.August, 27, 18, 11, 25, 0, time.UTC)

	tests := map[string]struct {
		givenTF                   string
		init                      func(*iam.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"happy path - role is returned by id": {
			givenTF: "valid_id.tf",
			init: func(m *iam.Mock) {
				m.On("GetRole", testutils.MockContext, iam.GetRoleRequest{
					ID:           12345,
					Actions:      true,
					GrantedRoles: true,
					Users:        true,
				}).Return(&iam.Role{
					RoleID:          int64(12345),
					RoleName:        "example-role",
					RoleDescription: "This is an example role.",
					CreatedBy:       "user@example.com",
					CreatedDate:     createdDate,
					ModifiedBy:      "admin@example.com",
					ModifiedDate:    modifiedDate,
					RoleType:        "custom",
					Actions: &iam.RoleAction{
						Delete: true,
						Edit:   true,
					},
				}, nil).Times(3)
			},
			expectedAttributes: map[string]string{
				"role_id":          "12345",
				"role_name":        "example-role",
				"role_description": "This is an example role.",
				"created_by":       "user@example.com",
				"created_date":     "2017-07-27T18:11:25Z",
				"modified_by":      "admin@example.com",
				"modified_date":    "2017-08-27T18:11:25Z",
				"type":             "custom",
				"actions.delete":   "true",
				"actions.edit":     "true",
			},
			expectError: nil,
		},
		"happy path - role is returned by name": {
			givenTF: "valid_name.tf",
			init: func(m *iam.Mock) {
				m.On("ListRoles", testutils.MockContext, iam.ListRolesRequest{}).Return([]iam.Role{{
					RoleID:          int64(12345),
					RoleName:        "example-role",
					RoleDescription: "This is an example role.",
					CreatedBy:       "user@example.com",
					CreatedDate:     createdDate,
					ModifiedBy:      "admin@example.com",
					ModifiedDate:    modifiedDate,
					RoleType:        "custom",
					Actions: &iam.RoleAction{
						Delete: true,
						Edit:   true,
					},
				},
				}, nil).Times(3)

				m.On("GetRole", testutils.MockContext, iam.GetRoleRequest{
					ID:           12345,
					Actions:      true,
					GrantedRoles: true,
					Users:        true,
				}).Return(&iam.Role{
					RoleID:          int64(12345),
					RoleName:        "example-role",
					RoleDescription: "This is an example role.",
					CreatedBy:       "user@example.com",
					CreatedDate:     createdDate,
					ModifiedBy:      "admin@example.com",
					ModifiedDate:    modifiedDate,
					RoleType:        "custom",
					Actions: &iam.RoleAction{
						Delete: true,
						Edit:   true,
					},
				}, nil).Times(3)
			},
			expectedAttributes: map[string]string{
				"role_id":          "12345",
				"role_name":        "example-role",
				"role_description": "This is an example role.",
				"created_by":       "user@example.com",
				"created_date":     "2017-07-27T18:11:25Z",
				"modified_by":      "admin@example.com",
				"modified_date":    "2017-08-27T18:11:25Z",
				"type":             "custom",
				"actions.delete":   "true",
				"actions.edit":     "true",
			},
			expectError: nil,
		},
		"happy path - role is returned, without dates": {
			givenTF: "valid_id.tf",
			init: func(m *iam.Mock) {
				m.On("GetRole", testutils.MockContext, iam.GetRoleRequest{
					ID:           12345,
					Actions:      true,
					GrantedRoles: true,
					Users:        true,
				}).Return(&iam.Role{
					RoleID:          int64(12345),
					RoleName:        "example-role",
					RoleDescription: "This is an example role.",
					CreatedBy:       "user@example.com",
					ModifiedBy:      "admin@example.com",
					RoleType:        "custom",
					Actions: &iam.RoleAction{
						Delete: true,
						Edit:   true,
					},
				}, nil).Times(3)
			},
			expectedAttributes: map[string]string{
				"role_id":          "12345",
				"role_name":        "example-role",
				"role_description": "This is an example role.",
				"created_by":       "user@example.com",
				"created_date":     "",
				"modified_by":      "admin@example.com",
				"modified_date":    "",
				"type":             "custom",
				"actions.delete":   "true",
				"actions.edit":     "true",
			},
			expectError: nil,
		},
		"error response from Get role endpoint": {
			givenTF: "valid_id.tf",
			init: func(m *iam.Mock) {
				m.On("GetRole", testutils.MockContext, iam.GetRoleRequest{
					ID:           12345,
					Actions:      true,
					GrantedRoles: true,
					Users:        true,
				}).Return(nil, fmt.Errorf("API error")).Once()
			},
			expectError: regexp.MustCompile("API error"),
		},
		"error response from List roles endpoint": {
			givenTF: "valid_name.tf",
			init: func(m *iam.Mock) {
				m.On("ListRoles", testutils.MockContext, iam.ListRolesRequest{}).Return(nil, fmt.Errorf("API error")).Once()
			},
			expectError: regexp.MustCompile("API error"),
		},
		"missing one of required arguments": {
			givenTF:     "missing_required_args.tf",
			expectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[role_name,role_id] is\srequired`),
		},
		"error - role not found by name": {
			givenTF: "valid_name.tf",
			init: func(m *iam.Mock) {
				m.On("ListRoles", testutils.MockContext, iam.ListRolesRequest{}).Return([]iam.Role{{
					RoleID:          int64(54321),
					RoleName:        "other-name",
					RoleDescription: "This is an example role.",
					CreatedBy:       "user@example.com",
					CreatedDate:     createdDate,
					ModifiedBy:      "admin@example.com",
					ModifiedDate:    modifiedDate,
					RoleType:        "custom",
					Actions: &iam.RoleAction{
						Delete: true,
						Edit:   true,
					},
				},
				}, nil).Once()
			},
			expectedAttributes: map[string]string{
				"role_id":          "12345",
				"role_name":        "example-role",
				"role_description": "This is an example role.",
				"created_by":       "user@example.com",
				"created_date":     "2017-07-27T18:11:25Z",
				"modified_by":      "admin@example.com",
				"modified_date":    "2017-08-27T18:11:25Z",
				"type":             "custom",
				"actions.delete":   "true",
				"actions.edit":     "true",
			},
			expectError: regexp.MustCompile("role with name 'example-role' not found"),
		},
		"error - more than one role found by name": {
			givenTF: "valid_name.tf",
			init: func(m *iam.Mock) {
				m.On("ListRoles", testutils.MockContext, iam.ListRolesRequest{}).Return([]iam.Role{{
					RoleID:          int64(12345),
					RoleName:        "example-role",
					RoleDescription: "This is an example role.",
					CreatedBy:       "user@example.com",
					CreatedDate:     createdDate,
					ModifiedBy:      "admin@example.com",
					ModifiedDate:    modifiedDate,
					RoleType:        "custom",
					Actions: &iam.RoleAction{
						Delete: true,
						Edit:   true,
					},
				},
					{
						RoleID:          int64(54321),
						RoleName:        "example-role",
						RoleDescription: "This is an 2nd example of role.",
						CreatedBy:       "user@example.com",
						CreatedDate:     createdDate,
						ModifiedBy:      "admin@example.com",
						ModifiedDate:    modifiedDate,
						RoleType:        "custom",
						Actions: &iam.RoleAction{
							Delete: true,
							Edit:   true,
						},
					},
				}, nil).Once()
			},
			expectedAttributes: map[string]string{
				"role_id":          "12345",
				"role_name":        "example-role",
				"role_description": "This is an example role.",
				"created_by":       "user@example.com",
				"created_date":     "2017-07-27T18:11:25Z",
				"modified_by":      "admin@example.com",
				"modified_date":    "2017-08-27T18:11:25Z",
				"type":             "custom",
				"actions.delete":   "true",
				"actions.edit":     "true",
			},
			expectError: regexp.MustCompile(`multiple roles with name 'example-role' found\. Specific roles IDs are '\[(?:12345\s54321|54321\s12345)\]'\. Please use 'role_id' to select the desired role`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}

			if tc.init != nil {
				tc.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range tc.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_role.test", k, v))
			}
			for _, v := range tc.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_iam_role.test", v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataRole/%s", tc.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: tc.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
