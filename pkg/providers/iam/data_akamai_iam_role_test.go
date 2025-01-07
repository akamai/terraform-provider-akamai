package iam

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
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
		"happy path - role is returned": {
			givenTF: "valid.tf",
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
		"happy path - role is returned, without dates": {
			givenTF: "valid.tf",
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
		"error response from API": {
			givenTF: "valid.tf",
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
		"missing required argument role_id": {
			givenTF:     "missing_role_id.tf",
			expectError: regexp.MustCompile(`The argument "role_id" is required, but no definition was found`),
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
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataRole/%s", tc.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: tc.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
