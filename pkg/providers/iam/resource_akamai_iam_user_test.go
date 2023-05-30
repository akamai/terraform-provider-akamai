package iam

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
	"github.com/stretchr/testify/assert"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceUser(t *testing.T) {
	basicUserInfo := iam.UserBasicInfo{
		FirstName:  "John",
		LastName:   "Smith",
		Email:      "jsmith@example.com",
		Phone:      "(111) 111-1111",
		TFAEnabled: false,
		Country:    "country",
	}
	extendedUserInfo := iam.UserBasicInfo{
		FirstName:         "John",
		LastName:          "Smith",
		Email:             "jsmith@example.com",
		Phone:             "(111) 111-1111",
		TimeZone:          "timezone",
		JobTitle:          "job title",
		TFAEnabled:        false,
		SecondaryEmail:    "secondary.email@example.com",
		MobilePhone:       "(222) 222-2222",
		Address:           "123 B Street",
		City:              "B-Town",
		State:             "state",
		ZipCode:           "zip",
		Country:           "country",
		ContactType:       "contact type",
		PreferredLanguage: "language",
		SessionTimeOut:    tools.IntPtr(2),
	}

	authGrantsCreate := []iam.AuthGrant{
		{
			GroupID:   0,
			GroupName: "group",
		},
	}
	authGrantsCreateRequest := []iam.AuthGrantRequest{
		{
			GroupID: 0,
		},
	}

	authGrantsSubgroupCreate := []iam.AuthGrant{
		{
			Subgroups: []iam.AuthGrant{
				{
					GroupID:   2,
					IsBlocked: false,
				},
				{
					GroupID:   1,
					IsBlocked: false,
				},
			},
		},
	}
	authGrantsSubgroupCreateRequest := []iam.AuthGrantRequest{
		{
			Subgroups: []iam.AuthGrantRequest{
				{
					GroupID:   2,
					IsBlocked: false,
				},
				{
					GroupID:   1,
					IsBlocked: false,
				},
			},
		},
	}

	authGrantsUpdate := []iam.AuthGrant{
		{
			GroupID:   1,
			GroupName: "other_group",
		},
	}

	authGrantsUpdateRequest := []iam.AuthGrantRequest{
		{
			GroupID: 1,
		},
	}

	notifications := iam.UserNotifications{
		Options: iam.UserNotificationOptions{
			Proactive: []string{},
			Upgrade:   []string{},
		},
	}
	id := "test_identity_id"

	checkUserAttributes := func(user iam.User) resource.TestCheckFunc {
		if user.SessionTimeOut == nil {
			user.SessionTimeOut = tools.IntPtr(0)
		}

		var authGrantsJSON string
		if len(user.AuthGrants) > 0 {
			marshalledAuthGrants, err := json.Marshal(user.AuthGrants)
			if err != nil {
				assert.NoError(t, err)
			}
			authGrantRequest := make([]iam.AuthGrantRequest, 0)
			err = json.Unmarshal(marshalledAuthGrants, &authGrantRequest)
			if err != nil {
				assert.NoError(t, err)
			}
			marshalledAuthGrants, err = json.Marshal(authGrantRequest)
			if err != nil {
				assert.NoError(t, err)
			}
			authGrantsJSON = string(marshalledAuthGrants)
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_iam_user.test", "id", "test_identity_id"),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "first_name", user.FirstName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "last_name", user.LastName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "email", strings.ToLower(user.Email)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "country", user.Country),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "phone", canonicalPhone(user.Phone)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "enable_tfa", fmt.Sprintf("%t", user.TFAEnabled)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "contact_type", user.ContactType),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "user_name", user.UserName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "job_title", user.JobTitle),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "time_zone", user.TimeZone),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "secondary_email", strings.ToLower(user.SecondaryEmail)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "mobile_phone", canonicalPhone(user.MobilePhone)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "address", user.Address),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "city", user.City),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "state", user.State),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "zip_code", user.ZipCode),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "preferred_language", user.PreferredLanguage),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "is_locked", fmt.Sprintf("%t", user.IsLocked)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "last_login", user.LastLoginDate),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "password_expired_after", user.PasswordExpiryDate),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "tfa_configured", fmt.Sprintf("%t", user.TFAConfigured)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "email_update_pending", fmt.Sprintf("%t", user.EmailUpdatePending)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "session_timeout", fmt.Sprintf("%d", *user.SessionTimeOut)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "auth_grants_json", authGrantsJSON),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "lock", fmt.Sprintf("%t", user.IsLocked)),
		)
	}

	userCreate := iam.User{
		UserBasicInfo:      basicUserInfo,
		IdentityID:         id,
		IsLocked:           false,
		LastLoginDate:      "last login",
		PasswordExpiryDate: "password expired after",
		TFAConfigured:      true,
		EmailUpdatePending: true,
		AuthGrants:         authGrantsCreate,
		Notifications:      notifications,
	}

	userSubgroupCreate := iam.User{
		UserBasicInfo:      basicUserInfo,
		IdentityID:         id,
		IsLocked:           false,
		LastLoginDate:      "last login",
		PasswordExpiryDate: "password expired after",
		TFAConfigured:      true,
		EmailUpdatePending: true,
		AuthGrants:         authGrantsSubgroupCreate,
		Notifications:      notifications,
	}

	userCreateRequest := iam.CreateUserRequest{
		UserBasicInfo: basicUserInfo,
		AuthGrants:    authGrantsCreateRequest,
		Notifications: notifications,
	}

	userSubgroupCreateRequest := iam.CreateUserRequest{
		UserBasicInfo: basicUserInfo,
		AuthGrants:    authGrantsSubgroupCreateRequest,
		Notifications: notifications,
	}

	userCreateLocked := iam.User{
		UserBasicInfo:      basicUserInfo,
		IdentityID:         id,
		IsLocked:           true,
		LastLoginDate:      "last login",
		PasswordExpiryDate: "password expired after",
		TFAConfigured:      true,
		EmailUpdatePending: true,
		AuthGrants:         authGrantsCreate,
		Notifications:      notifications,
	}

	userCreateLockedRequest := iam.CreateUserRequest{
		UserBasicInfo: basicUserInfo,
		AuthGrants:    authGrantsCreateRequest,
		Notifications: notifications,
	}

	userUpdateInfo := iam.User{
		UserBasicInfo:      extendedUserInfo,
		IdentityID:         id,
		IsLocked:           false,
		LastLoginDate:      "last login",
		PasswordExpiryDate: "password expired after",
		TFAConfigured:      true,
		EmailUpdatePending: true,
		AuthGrants:         authGrantsCreate,
		Notifications:      notifications,
	}
	userUpdateGrants := iam.User{
		UserBasicInfo:      basicUserInfo,
		IdentityID:         id,
		IsLocked:           false,
		LastLoginDate:      "last login",
		PasswordExpiryDate: "password expired after",
		TFAConfigured:      true,
		EmailUpdatePending: true,
		AuthGrants:         authGrantsUpdate,
		Notifications:      notifications,
	}

	userUpdateGrantsRequest := iam.CreateUserRequest{
		UserBasicInfo: basicUserInfo,
		AuthGrants:    authGrantsUpdateRequest,
		Notifications: notifications,
	}
	authGrantsCreateWithIgnoredFields := []iam.AuthGrantRequest{
		{
			GroupID:   1,
			IsBlocked: false,
		},
	}
	authGrantsCreateWithIgnoredFieldsResponse := []iam.AuthGrant{
		{
			GroupID:         1,
			GroupName:       "group",
			IsBlocked:       false,
			RoleDescription: "desc from server",
			RoleID:          nil,
			RoleName:        "role name from server",
			Subgroups:       nil,
		},
	}

	userCreateWithIgnoredFields := userCreate
	userCreateWithIgnoredFieldsRequest := userCreateRequest
	userCreateWithIgnoredFieldsRequest.AuthGrants = authGrantsCreateWithIgnoredFields
	userCreateWithIgnoredFieldsResponse := userCreate
	userCreateWithIgnoredFieldsResponse.AuthGrants = authGrantsCreateWithIgnoredFieldsResponse

	tests := map[string]struct {
		init  func(*iam.Mock)
		steps []resource.TestStep
	}{
		"basic": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userCreate, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
			},
		},
		"basic lock": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateLockedRequest, userCreateLocked, true, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreateLocked, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userCreateLocked, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic_lock.tf"),
					Check:  checkUserAttributes(userCreateLocked),
				},
			},
		},
		"basic error create": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, fmt.Errorf("error create"), nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					ExpectError: regexp.MustCompile("failed to create user: error create"),
				},
			},
		},
		"basic no diff no update": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userCreate, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
			},
		},
		"update user info": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update basic info
				expectResourceIAMUserInfoUpdatePhase(m, userUpdateInfo.IdentityID, userUpdateInfo.UserBasicInfo, nil).Once()
				expectResourceIAMUserReadPhase(m, userUpdateInfo, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateInfo, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/update_user_info.tf"),
					Check:  checkUserAttributes(userUpdateInfo),
				},
			},
		},
		"update user info - lock - unlock": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update lock
				expectResourceIAMUserReadPhase(m, userCreateLocked, nil).Once()

				// plan
				expectResourceIAMUserReadPhase(m, userCreateLocked, nil).Once()
				// update lock
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userCreate, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic_lock.tf"),
					Check:  checkUserAttributes(userCreateLocked),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
			},
		},
		"update user info - error": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update basic info
				expectResourceIAMUserInfoUpdatePhase(m, userUpdateInfo.IdentityID, userUpdateInfo.UserBasicInfo, fmt.Errorf("error updating")).Once()

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateInfo, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceUserLifecycle/update_user_info.tf"),
					ExpectError: regexp.MustCompile("failed to update user: error updating"),
				},
			},
		},
		"update user auth grants": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update basic info
				expectResourceIAMUserAuthGrantsUpdatePhase(m, userUpdateGrants.IdentityID, userUpdateGrantsRequest.AuthGrants, userUpdateGrants.AuthGrants, nil).Once()
				expectResourceIAMUserReadPhase(m, userUpdateGrants, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateGrants, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/update_auth_grants.tf"),
					Check:  checkUserAttributes(userUpdateGrants),
				},
			},
		},
		"update swap user auth grants subgroups": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userSubgroupCreateRequest, userSubgroupCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userSubgroupCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userSubgroupCreate, nil).Times(2)

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateGrants, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic_grants.tf"),
					Check:  checkUserAttributes(userSubgroupCreate),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic_grants_swap.tf"),
					Check:  checkUserAttributes(userSubgroupCreate),
				},
			},
		},
		"update user auth grants with redundant fields": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateWithIgnoredFieldsRequest, userCreateWithIgnoredFields, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreateWithIgnoredFieldsResponse, nil).Once()
				expectResourceIAMUserReadPhase(m, userCreateWithIgnoredFieldsResponse, nil).Once()

				// plan
				expectResourceIAMUserReadPhase(m, userCreateWithIgnoredFieldsResponse, nil).Once()
				expectResourceIAMUserReadPhase(m, userCreateWithIgnoredFieldsResponse, nil).Once()

				// delete
				expectResourceIAMUserDeletePhase(m, userCreateWithIgnoredFieldsResponse, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_auth_grants.tf"),
					Check:  checkUserAttributes(userCreateWithIgnoredFieldsResponse),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_auth_grants.tf"),
					Check:  checkUserAttributes(userCreateWithIgnoredFieldsResponse),
				},
			},
		},
		"update user auth grants - an error": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update basic info
				expectResourceIAMUserAuthGrantsUpdatePhase(m, userUpdateGrants.IdentityID, userUpdateGrantsRequest.AuthGrants, userUpdateGrants.AuthGrants, fmt.Errorf("error update user auth grants")).Once()

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateGrants, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceUserLifecycle/update_auth_grants.tf"),
					ExpectError: regexp.MustCompile("failed to update user AuthGrants: error update user auth grants"),
				},
			},
		},
		"basic import": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// import
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(1)

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateGrants, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_iam_user.test",
					ImportStateVerify: true,
				},
			},
		},
		"auth_grants_json should not panic when supplied interpolated string with unknown value": {
			init: func(m *iam.Mock) {
			},
			steps: []resource.TestStep{
				{
					Config:             loadFixtureString("./testdata/TestResourceUserLifecycle/auth_grants_interpolated.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"error updating email": {
			init: func(m *iam.Mock) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateRequest, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()

				// delete
				expectResourceIAMUserDeletePhase(m, userUpdateGrants, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					Check:  checkUserAttributes(userCreate),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceUserLifecycle/update_email.tf"),
					ExpectError: regexp.MustCompile("cannot change email address"),
				},
			},
		},
		"error creating user: invalid auth grants": {
			init: func(m *iam.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceUserLifecycle/invalid_auth_grants.tf"),
					ExpectError: regexp.MustCompile("auth_grants_json is not valid"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps:             test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// create
func expectResourceIAMUserCreatePhase(m *iam.Mock, request iam.CreateUserRequest, response iam.User, lock bool, creationError, lockError error) {
	onCreation := m.On("CreateUser", mock.Anything, iam.CreateUserRequest{
		UserBasicInfo: request.UserBasicInfo,
		AuthGrants:    request.AuthGrants,
		SendEmail:     true,
		Notifications: request.Notifications,
	})
	if creationError != nil {
		onCreation.Return(nil, creationError).Once()
		return
	}
	onCreation.Return(&response, nil).Once()

	if lock {
		expectToggleLock(m, response.IdentityID, true, lockError).Once()
		if lockError != nil {
			return
		}
	}
}

func expectToggleLock(m *iam.Mock, identityID string, lock bool, err error) *mock.Call {
	if lock {
		return m.On("LockUser", mock.Anything, iam.LockUserRequest{IdentityID: identityID}).Return(err)
	}
	return m.On("UnlockUser", mock.Anything, iam.UnlockUserRequest{IdentityID: identityID}).Return(err)
}

// read
func expectResourceIAMUserReadPhase(m *iam.Mock, user iam.User, anError error) *mock.Call {
	on := m.On("GetUser", mock.Anything, iam.GetUserRequest{
		IdentityID: user.IdentityID,
		AuthGrants: true,
	})
	if anError != nil {
		return on.Return(nil, anError).Once()
	}
	return on.Return(&user, nil)
}

// update user info
func expectResourceIAMUserInfoUpdatePhase(m *iam.Mock, id string, basicUserInfo iam.UserBasicInfo, anError error) *mock.Call {
	on := m.On("UpdateUserInfo", mock.Anything, iam.UpdateUserInfoRequest{
		IdentityID: id,
		User:       basicUserInfo,
	})
	if anError != nil {
		return on.Return(nil, anError).Once()
	}
	return on.Return(&basicUserInfo, nil)
}

// update auth grants
func expectResourceIAMUserAuthGrantsUpdatePhase(m *iam.Mock, id string, authGrantsReqest []iam.AuthGrantRequest, authGrants []iam.AuthGrant, anError error) *mock.Call {
	on := m.On("UpdateUserAuthGrants", mock.Anything, iam.UpdateUserAuthGrantsRequest{
		IdentityID: id,
		AuthGrants: authGrantsReqest,
	})
	if anError != nil {
		return on.Return(nil, anError).Once()
	}
	return on.Return(authGrants, nil)
}

// delete
func expectResourceIAMUserDeletePhase(m *iam.Mock, user iam.User, anError error) *mock.Call {
	on := m.On("RemoveUser", mock.Anything, iam.RemoveUserRequest{IdentityID: user.IdentityID})
	if anError != nil {
		return on.Return(anError).Once()
	}
	return on.Return(nil)
}
