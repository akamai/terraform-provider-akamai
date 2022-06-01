package iam

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/stretchr/testify/assert"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
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
	authGrantsUpdate := []iam.AuthGrant{
		{
			GroupID:   1,
			GroupName: "other_group",
		},
	}

	notifications := iam.UserNotifications{
		Options: iam.UserNotificationOptions{
			Proactive: []string{},
			Upgrade:   []string{},
		},
	}
	id := "test_identity_id"

	checkUserAttributes := func(User iam.User) resource.TestCheckFunc {
		if User.SessionTimeOut == nil {
			User.SessionTimeOut = tools.IntPtr(0)
		}

		var authGrantsJSON string
		if len(User.AuthGrants) > 0 {
			asd, err := json.Marshal(User.AuthGrants)
			if err != nil {
				assert.NoError(t, err)
			}
			authGrantsJSON = string(asd)
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_iam_user.test", "id", "test_identity_id"),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "first_name", User.FirstName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "last_name", User.LastName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "email", strings.ToLower(User.Email)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "country", User.Country),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "phone", canonicalPhone(User.Phone)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "enable_tfa", fmt.Sprintf("%t", User.TFAEnabled)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "contact_type", User.ContactType),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "user_name", User.UserName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "job_title", User.JobTitle),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "time_zone", User.TimeZone),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "secondary_email", strings.ToLower(User.SecondaryEmail)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "mobile_phone", canonicalPhone(User.MobilePhone)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "address", User.Address),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "city", User.City),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "state", User.State),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "zip_code", User.ZipCode),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "preferred_language", User.PreferredLanguage),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "is_locked", fmt.Sprintf("%t", User.IsLocked)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "last_login", User.LastLoginDate),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "password_expired_after", User.PasswordExpiryDate),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "tfa_configured", fmt.Sprintf("%t", User.TFAConfigured)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "email_update_pending", fmt.Sprintf("%t", User.EmailUpdatePending)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "session_timeout", fmt.Sprintf("%d", *User.SessionTimeOut)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "auth_grants_json", authGrantsJSON),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "lock", fmt.Sprintf("%t", User.IsLocked)),
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

	tests := map[string]struct {
		init  func(*mockiam)
		steps []resource.TestStep
	}{
		"basic": {
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreateLocked, true, nil, nil)
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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, fmt.Errorf("error create"), nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceUserLifecycle/create_basic.tf"),
					ExpectError: regexp.MustCompile("failed to create user: error create"),
				},
			},
		},
		"basic no diff no update": {
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update basic info
				expectResourceIAMUserAuthGrantsUpdatePhase(m, userUpdateGrants.IdentityID, userUpdateGrants.AuthGrants, nil).Once()
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
		"update user auth grants - an error": {
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
				expectResourceIAMUserReadPhase(m, userCreate, nil).Times(2)

				// plan
				expectResourceIAMUserReadPhase(m, userCreate, nil).Once()
				// update basic info
				expectResourceIAMUserAuthGrantsUpdatePhase(m, userUpdateGrants.IdentityID, userUpdateGrants.AuthGrants, fmt.Errorf("error update user auth grants")).Once()

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
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
		"error updating email": {
			init: func(m *mockiam) {
				// create
				expectResourceIAMUserCreatePhase(m, userCreate, false, nil, nil)
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
			init: func(m *mockiam) {},
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
			client := &mockiam{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// create
func expectResourceIAMUserCreatePhase(m *mockiam, user iam.User, lock bool, creationError, lockError error) {
	onCreation := m.On("CreateUser", mock.Anything, iam.CreateUserRequest{
		User:          user.UserBasicInfo,
		AuthGrants:    user.AuthGrants,
		SendEmail:     true,
		Notifications: user.Notifications,
	})
	if creationError != nil {
		onCreation.Return(nil, creationError).Once()
		return
	}
	onCreation.Return(&user, nil).Once()

	if lock {
		expectToggleLock(m, user.IdentityID, true, lockError).Once()
		if lockError != nil {
			return
		}
	}
}

func expectToggleLock(m *mockiam, identityID string, lock bool, err error) *mock.Call {
	if lock {
		return m.On("LockUser", mock.Anything, iam.LockUserRequest{IdentityID: identityID}).Return(err)
	}
	return m.On("UnlockUser", mock.Anything, iam.UnlockUserRequest{IdentityID: identityID}).Return(err)
}

// read
func expectResourceIAMUserReadPhase(m *mockiam, user iam.User, anError error) *mock.Call {
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
func expectResourceIAMUserInfoUpdatePhase(m *mockiam, id string, basicUserInfo iam.UserBasicInfo, anError error) *mock.Call {
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
func expectResourceIAMUserAuthGrantsUpdatePhase(m *mockiam, id string, authGrants []iam.AuthGrant, anError error) *mock.Call {
	on := m.On("UpdateUserAuthGrants", mock.Anything, iam.UpdateUserAuthGrantsRequest{
		IdentityID: id,
		AuthGrants: authGrants,
	})
	if anError != nil {
		return on.Return(nil, anError).Once()
	}
	return on.Return(authGrants, nil)
}

// delete
func expectResourceIAMUserDeletePhase(m *mockiam, user iam.User, anError error) *mock.Call {
	on := m.On("RemoveUser", mock.Anything, iam.RemoveUserRequest{IdentityID: user.IdentityID})
	if anError != nil {
		return on.Return(anError).Once()
	}
	return on.Return(nil)
}
