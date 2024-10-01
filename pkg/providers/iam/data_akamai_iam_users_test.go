package iam

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

var (
	basicTestDataForUsers = []iam.UserListItem{{
		FirstName:     "firstName1",
		LastName:      "lastName1",
		UserName:      "user1",
		Email:         "user1@example.com",
		TFAEnabled:    false,
		IdentityID:    "A-B-123456",
		IsLocked:      false,
		LastLoginDate: time.Date(2017, time.July, 27, 18, 11, 25, 0, time.UTC),
		TFAConfigured: false,
		AccountID:     "A-CCT3456",
		Actions: &iam.UserActions{
			APIClient:        true,
			Delete:           true,
			Edit:             true,
			IsCloneable:      true,
			ResetPassword:    true,
			ThirdPartyAccess: true,
			CanEditTFA:       false,
			CanEditMFA:       false,
			CanEditNone:      false,
			EditProfile:      true,
		},
		AuthGrants: []iam.AuthGrant{
			{
				GroupID:         12345,
				GroupName:       "Internet Company",
				IsBlocked:       false,
				RoleDescription: "Role giving admin access.",
				RoleID:          ptr.To(12),
				RoleName:        "admin",
				Subgroups: []iam.AuthGrant{{
					GroupID:   1000,
					GroupName: "sub group",
					RoleID:    ptr.To(1000),
					RoleName:  "sub group role",
				}},
			},
			{
				RoleID:   ptr.To(100002),
				RoleName: "Admin for Account Roles",
			},
		},
		AdditionalAuthentication:           "TFA",
		AdditionalAuthenticationConfigured: false,
	},
		{
			FirstName:     "firstName2",
			LastName:      "lastName2",
			UserName:      "user2",
			Email:         "user2@example.com",
			TFAEnabled:    false,
			IdentityID:    "B-B-123456",
			IsLocked:      false,
			LastLoginDate: time.Time{},
			TFAConfigured: false,
			AccountID:     "B-CCT3456",
			Actions: &iam.UserActions{
				APIClient:        true,
				Delete:           true,
				Edit:             true,
				IsCloneable:      true,
				ResetPassword:    true,
				ThirdPartyAccess: true,
				CanEditTFA:       false,
				CanEditMFA:       false,
				CanEditNone:      false,
				EditProfile:      true,
			},
			AuthGrants: []iam.AuthGrant{
				{
					RoleID:   ptr.To(100002),
					RoleName: "Admin for Account Roles",
				},
			},
			AdditionalAuthentication:           "TFA",
			AdditionalAuthenticationConfigured: false,
		}}
)

func TestDataUsers(t *testing.T) {
	tests := map[string]struct {
		configPath    string
		init          func(*testing.T, *iam.Mock, []iam.UserListItem, *int64)
		mockData      []iam.UserListItem
		groupID       *int64
		expectedError *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataUsers/default.tf",
			init: func(t *testing.T, m *iam.Mock, mockData []iam.UserListItem, groupID *int64) {
				expectListUsers(m, mockData, groupID, 3)
			},
			mockData: basicTestDataForUsers,
		},
		"happy path - no users": {
			configPath: "testdata/TestDataUsers/default.tf",
			init: func(t *testing.T, m *iam.Mock, mockData []iam.UserListItem, groupID *int64) {
				expectListUsers(m, mockData, groupID, 3)
			},
			mockData: []iam.UserListItem{},
		},
		"happy path - groupID search": {
			groupID:    ptr.To(int64(12345)),
			configPath: "testdata/TestDataUsers/groupIDSearch.tf",
			init: func(t *testing.T, m *iam.Mock, mockData []iam.UserListItem, groupID *int64) {
				expectListUsers(m, mockData, groupID, 3)
			},
			mockData: basicTestDataForUsers,
		},
		"error - list user fails": {
			configPath: "testdata/TestDataUsers/default.tf",
			init: func(t *testing.T, m *iam.Mock, mockData []iam.UserListItem, groupID *int64) {
				listUsersReq := iam.ListUsersRequest{GroupID: groupID, AuthGrants: true, Actions: true}
				m.On("ListUsers", mock.Anything, listUsersReq).Return(nil, errors.New("test error"))
			},
			expectedError: regexp.MustCompile("test error"),
			mockData:      basicTestDataForUsers,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(t, client, tc.mockData, tc.groupID)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.configPath),
							Check:       checkUsersAttrs(tc.groupID, len(tc.mockData) == 0),
							ExpectError: tc.expectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkUsersAttrs(groupID *int64, emptyReturn bool) resource.TestCheckFunc {
	name := "data.akamai_iam_users.test"

	if emptyReturn {
		checksFuncs := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(name, "users.#", "0"),
		}
		if groupID != nil {
			checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "group_id", "12345"))
		} else {
			checksFuncs = append(checksFuncs, resource.TestCheckNoResourceAttr(name, "group_id"))
		}
		return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
	}

	checksFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "users.#", "2"),
		resource.TestCheckResourceAttr(name, "users.0.account_id", "A-CCT3456"),
		resource.TestCheckResourceAttr(name, "users.0.additional_authentication", "TFA"),
		resource.TestCheckResourceAttr(name, "users.0.additional_authentication_configured", "false"),
		resource.TestCheckResourceAttr(name, "users.0.email", "user1@example.com"),
		resource.TestCheckResourceAttr(name, "users.0.first_name", "firstName1"),
		resource.TestCheckResourceAttr(name, "users.0.is_locked", "false"),
		resource.TestCheckResourceAttr(name, "users.0.last_login_date", "2017-07-27T18:11:25Z"),
		resource.TestCheckResourceAttr(name, "users.0.last_name", "lastName1"),
		resource.TestCheckResourceAttr(name, "users.0.tfa_configured", "false"),
		resource.TestCheckResourceAttr(name, "users.0.tfa_enabled", "false"),
		resource.TestCheckResourceAttr(name, "users.0.ui_identity_id", "A-B-123456"),
		resource.TestCheckResourceAttr(name, "users.0.ui_user_name", "user1"),
		resource.TestCheckResourceAttr(name, "users.0.actions.delete", "true"),
		resource.TestCheckResourceAttr(name, "users.0.actions.api_client", "true"),
		resource.TestCheckResourceAttr(name, "users.0.actions.can_edit_mfa", "false"),
		resource.TestCheckResourceAttr(name, "users.0.actions.can_edit_none", "false"),
		resource.TestCheckResourceAttr(name, "users.0.actions.can_edit_tfa", "false"),
		resource.TestCheckResourceAttr(name, "users.0.actions.edit", "true"),
		resource.TestCheckResourceAttr(name, "users.0.actions.edit_profile", "true"),
		resource.TestCheckResourceAttr(name, "users.0.actions.is_cloneable", "true"),
		resource.TestCheckResourceAttr(name, "users.0.actions.reset_password", "true"),
		resource.TestCheckResourceAttr(name, "users.0.actions.third_party_access", "true"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.#", "2"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.group_id", "12345"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.group_name", "Internet Company"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.is_blocked", "false"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.role_description", "Role giving admin access."),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.role_id", "12"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.role_name", "admin"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.0.sub_groups.#", "1"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.group_id", "0"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.group_name", ""),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.is_blocked", "false"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.role_description", ""),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.role_id", "100002"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.role_name", "Admin for Account Roles"),
		resource.TestCheckResourceAttr(name, "users.0.auth_grants.1.sub_groups.#", "0"),

		resource.TestCheckResourceAttr(name, "users.1.account_id", "B-CCT3456"),
		resource.TestCheckResourceAttr(name, "users.1.additional_authentication", "TFA"),
		resource.TestCheckResourceAttr(name, "users.1.additional_authentication_configured", "false"),
		resource.TestCheckResourceAttr(name, "users.1.email", "user2@example.com"),
		resource.TestCheckResourceAttr(name, "users.1.first_name", "firstName2"),
		resource.TestCheckResourceAttr(name, "users.1.is_locked", "false"),
		resource.TestCheckResourceAttr(name, "users.1.last_login_date", ""),
		resource.TestCheckResourceAttr(name, "users.1.last_name", "lastName2"),
		resource.TestCheckResourceAttr(name, "users.1.tfa_configured", "false"),
		resource.TestCheckResourceAttr(name, "users.1.tfa_enabled", "false"),
		resource.TestCheckResourceAttr(name, "users.1.ui_identity_id", "B-B-123456"),
		resource.TestCheckResourceAttr(name, "users.1.ui_user_name", "user2"),
		resource.TestCheckResourceAttr(name, "users.1.actions.delete", "true"),
		resource.TestCheckResourceAttr(name, "users.1.actions.api_client", "true"),
		resource.TestCheckResourceAttr(name, "users.1.actions.can_edit_mfa", "false"),
		resource.TestCheckResourceAttr(name, "users.1.actions.can_edit_none", "false"),
		resource.TestCheckResourceAttr(name, "users.1.actions.can_edit_tfa", "false"),
		resource.TestCheckResourceAttr(name, "users.1.actions.edit", "true"),
		resource.TestCheckResourceAttr(name, "users.1.actions.edit_profile", "true"),
		resource.TestCheckResourceAttr(name, "users.1.actions.is_cloneable", "true"),
		resource.TestCheckResourceAttr(name, "users.1.actions.reset_password", "true"),
		resource.TestCheckResourceAttr(name, "users.1.actions.third_party_access", "true"),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.#", "1"),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.group_id", "0"),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.group_name", ""),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.is_blocked", "false"),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.role_description", ""),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.role_id", "100002"),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.role_name", "Admin for Account Roles"),
		resource.TestCheckResourceAttr(name, "users.1.auth_grants.0.sub_groups.#", "0"),
	}
	if groupID != nil {
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "group_id", "12345"))
	} else {
		checksFuncs = append(checksFuncs, resource.TestCheckNoResourceAttr(name, "group_id"))
	}

	return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
}

func expectListUsers(client *iam.Mock, mockData []iam.UserListItem, groupID *int64, timesToRun int) {
	listUsersReq := iam.ListUsersRequest{
		GroupID:    groupID,
		AuthGrants: true,
		Actions:    true,
	}

	client.On("ListUsers", mock.Anything, listUsersReq).Return(mockData, nil).Times(timesToRun)
}
