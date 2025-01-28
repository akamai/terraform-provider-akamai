package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataPropertyUsers(t *testing.T) {
	userJohn := iam.UsersForProperty{
		FirstName:    "John",
		LastName:     "Smith",
		IsBlocked:    false,
		UIIdentityID: "B-C-IP9IND",
		UIUserName:   "josmith",
	}

	userJane := iam.UsersForProperty{
		FirstName:    "Jane",
		LastName:     "Smith",
		IsBlocked:    false,
		UIIdentityID: "B-C-AB1CDE",
		UIUserName:   "jasmith",
	}

	blockedUserJudy := iam.UsersForProperty{
		FirstName:    "Judy",
		LastName:     "Smith",
		IsBlocked:    true,
		UIIdentityID: "B-C-FG2HIJ",
		UIUserName:   "jusmith",
	}

	mockListUsersForProperty := func(iamMock *iam.Mock, req iam.ListUsersForPropertyRequest,
		users ...iam.UsersForProperty) {
		iamMock.On("ListUsersForProperty", testutils.MockContext, req).
			Return((iam.ListUsersForPropertyResponse)(users), nil).Times(3)
	}

	tests := map[string]struct {
		configPath string
		init       func(*iam.Mock)
		check      resource.TestCheckFunc
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataPropertyUsers/basic.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
				}
				mockListUsersForProperty(iamMock, req, userJohn, userJane, blockedUserJudy)
			},
			check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "asset_id", "12345"),
				resource.TestCheckNoResourceAttr("data.akamai_iam_property_users.test", "user_type"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.#", "3"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.first_name", "John"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.last_name", "Smith"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.1.is_blocked", "false"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.2.ui_identity_id", "B-C-FG2HIJ"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.2.ui_user_name", "jusmith")),
		},
		"handling aid prefix": {
			configPath: "testdata/TestDataPropertyUsers/aid_prefix.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
				}
				mockListUsersForProperty(iamMock, req, userJohn, userJane, blockedUserJudy)
			},
			check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "asset_id", "aid_12345"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.#", "3"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.first_name", "John")),
		},
		"user filtering - two users assigned": {
			configPath: "testdata/TestDataPropertyUsers/user_type_assigned.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
					UserType:   "assigned",
				}
				mockListUsersForProperty(iamMock, req, userJohn, userJane)
			},
			check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "asset_id", "12345"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "user_type", "assigned"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.#", "2"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.first_name", "John"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.is_blocked", "false"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.1.first_name", "Jane"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.1.is_blocked", "false")),
		},
		"user filtering - one user blocked": {
			configPath: "testdata/TestDataPropertyUsers/user_type_blocked.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
					UserType:   "blocked",
				}
				mockListUsersForProperty(iamMock, req, blockedUserJudy)
			},
			check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "asset_id", "12345"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "user_type", "blocked"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.first_name", "Judy"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.is_blocked", "true")),
		},
		"user filtering - all users returned for all": {
			configPath: "testdata/TestDataPropertyUsers/user_type_all.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
					UserType:   "all",
				}
				mockListUsersForProperty(iamMock, req, userJohn, userJane, blockedUserJudy)
			},
			check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "asset_id", "12345"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "user_type", "all"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.#", "3"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.0.first_name", "John"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.2.first_name", "Judy")),
		},
		"no users returned": {
			configPath: "testdata/TestDataPropertyUsers/basic.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
				}
				mockListUsersForProperty(iamMock, req)
			},
			check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "asset_id", "12345"),
				resource.TestCheckNoResourceAttr("data.akamai_iam_property_users.test", "user_type"),
				resource.TestCheckResourceAttr("data.akamai_iam_property_users.test", "users.#", "0")),
		},
		"missing property id": {
			configPath: "testdata/TestDataPropertyUsers/missing_asset_id.tf",
			error:      regexp.MustCompile(`The argument "asset_id" is required, but no definition was found`),
		},
		"empty property id": {
			configPath: "testdata/TestDataPropertyUsers/empty_asset_id.tf",
			error:      regexp.MustCompile(`Attribute asset_id must be a number with the optional "aid_" prefix`),
		},
		"property id is not a number": {
			configPath: "testdata/TestDataPropertyUsers/nan_asset_id.tf",
			error:      regexp.MustCompile(`Attribute asset_id must be a number with the optional "aid_" prefix`),
		},
		"property id has invalid prefix": {
			configPath: "testdata/TestDataPropertyUsers/bad_prefix_asset_id.tf",
			error:      regexp.MustCompile(`Attribute asset_id must be a number with the optional "aid_" prefix`),
		},
		"property id has trailing text": {
			configPath: "testdata/TestDataPropertyUsers/trailing_text_asset_id.tf",
			error:      regexp.MustCompile(`Attribute asset_id must be a number with the optional "aid_" prefix`),
		},
		"bad user type": {
			configPath: "testdata/TestDataPropertyUsers/bad_user_type.tf",
			error:      regexp.MustCompile(`Attribute user_type value must be one of.+all.+blocked.+assigned`),
		},
		"edgegrid error": {
			configPath: "testdata/TestDataPropertyUsers/basic.tf",
			init: func(iamMock *iam.Mock) {
				req := iam.ListUsersForPropertyRequest{
					PropertyID: 12345,
				}
				iamMock.On("ListUsersForProperty", testutils.MockContext, req).
					Return(nil, fmt.Errorf("list users failed"))
			},
			error: regexp.MustCompile("list users failed"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			iamMock := iam.Mock{}
			if tc.init != nil {
				tc.init(&iamMock)
			}

			useClient(&iamMock, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.configPath),
							Check:       tc.check,
							ExpectError: tc.error,
						},
					},
				})
			})
			iamMock.AssertExpectations(t)
		})
	}
}
