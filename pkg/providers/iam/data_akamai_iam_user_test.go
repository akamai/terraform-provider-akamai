package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	testDataForUser struct {
		uiIdentityID                       string
		accountID                          string
		additionalAuthentication           string
		additionalAuthenticationConfigured bool
		address                            string
		city                               string
		contactType                        string
		country                            string
		email                              string
		emailUpdatePending                 bool
		firstName                          string
		isLocked                           bool
		jobTitle                           string
		lastLoginDate                      time.Time
		lastName                           string
		mobilePhone                        string
		passwordExpiryDate                 time.Time
		phone                              string
		preferredLanguage                  string
		secondaryEmail                     string
		sessionTimeout                     int64
		state                              string
		tfaConfigured                      bool
		tfaEnabled                         bool
		timeZone                           string
		uiUserName                         string
		zipCode                            string
		actions                            actionsData
		authGrants                         []authGrantData
		notifications                      notificationData
	}

	actionsData struct {
		delete           bool
		apiClient        bool
		edit             bool
		isCloneable      bool
		resetPassword    bool
		thirdPartyAccess bool
	}

	authGrantData struct {
		groupID         int64
		groupName       string
		isBlocked       bool
		roleDescription string
		roleID          int64
		roleName        string
		subgroup        []*authGrantData
	}
	notificationData struct {
		options                  optionsData
		enableEmailNotifications bool
	}

	optionsData struct {
		apiClientCredentialExpiryNotification bool
		newUserNotification                   bool
		passwordExpiry                        bool
		proactive                             []string
		upgrade                               []string
	}
)

var (
	basicUserTestData = testDataForUser{
		uiIdentityID:                       "asd-12345",
		uiUserName:                         "JoeDoeOh",
		accountID:                          "acc-12345",
		additionalAuthentication:           "NONE",
		additionalAuthenticationConfigured: true,
		address:                            "test address 12",
		city:                               "test city",
		contactType:                        "test contact type",
		country:                            "test country",
		email:                              "email@test.com",
		emailUpdatePending:                 true,
		firstName:                          "Joe",
		isLocked:                           true,
		jobTitle:                           "Phd",
		lastLoginDate:                      time.Date(2021, 1, 11, 7, 45, 18, 000, time.UTC),
		lastName:                           "Doe",
		mobilePhone:                        "123-456-789",
		passwordExpiryDate:                 time.Date(2025, 1, 11, 7, 45, 18, 000, time.UTC),
		phone:                              "987-654-321",
		preferredLanguage:                  "English",
		secondaryEmail:                     "seccondEmail@test.com",
		sessionTimeout:                     1000,
		state:                              "CN",
		tfaEnabled:                         true,
		tfaConfigured:                      true,
		timeZone:                           "UTC",
		zipCode:                            "12-345",
		actions:                            basicActionTestData,
		notifications:                      basicNotificationTesData,
		authGrants:                         append([]authGrantData{}, basicAuthGrantData),
	}
	basicUserTestDataMaxGroups = testDataForUser{
		uiIdentityID:                       "asd-12345",
		uiUserName:                         "JoeDoeOh",
		accountID:                          "acc-12345",
		additionalAuthentication:           "NONE",
		additionalAuthenticationConfigured: true,
		address:                            "test address 12",
		city:                               "test city",
		contactType:                        "test contact type",
		country:                            "test country",
		email:                              "email@test.com",
		emailUpdatePending:                 true,
		firstName:                          "Joe",
		isLocked:                           true,
		jobTitle:                           "Phd",
		lastLoginDate:                      time.Date(2021, 1, 11, 7, 45, 18, 000, time.UTC),
		lastName:                           "Doe",
		mobilePhone:                        "123-456-789",
		passwordExpiryDate:                 time.Date(2025, 1, 11, 7, 45, 18, 000, time.UTC),
		phone:                              "987-654-321",
		preferredLanguage:                  "English",
		secondaryEmail:                     "seccondEmail@test.com",
		sessionTimeout:                     1000,
		state:                              "CN",
		tfaEnabled:                         true,
		tfaConfigured:                      true,
		timeZone:                           "UTC",
		zipCode:                            "12-345",
		actions:                            basicActionTestData,
		notifications:                      basicNotificationTesData,
		authGrants:                         append([]authGrantData{}, basicAuthGrantDataMaxSubgroup),
	}

	basicActionTestData = actionsData{
		delete:           true,
		apiClient:        true,
		edit:             true,
		isCloneable:      true,
		resetPassword:    true,
		thirdPartyAccess: true,
	}

	basicNotificationTesData = notificationData{
		enableEmailNotifications: true,
		options:                  basicOptionsTestData,
	}

	basicOptionsTestData = optionsData{
		apiClientCredentialExpiryNotification: true,
		newUserNotification:                   true,
		passwordExpiry:                        true,
		proactive:                             []string{"EdgeScape"},
		upgrade:                               []string{"NetStorage"},
	}

	basicAuthGrantData = authGrantData{
		roleDescription: "testDesc",
		roleName:        "admin",
		isBlocked:       false,
		roleID:          1234,
		groupID:         1234,
		groupName:       "TestName",
	}
	basicAuthGrantDataMaxSubgroup = authGrantData{
		roleDescription: "testDesc",
		roleName:        "admin",
		isBlocked:       false,
		roleID:          1234,
		groupID:         1234,
		groupName:       "TestName",
		subgroup:        generateMaxDepthSubGroupsAuthGrantData(maxSupportedGroupNesting),
	}
)

func TestDataUser(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *iam.Mock, testDataForUser)
		mockData   testDataForUser
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataUser/default.tf",
			init: func(t *testing.T, m *iam.Mock, testData testDataForUser) {
				expectGetUser(t, m, testData, 3)
			},
			mockData: basicUserTestData,
		},
		"happy path - max amount of sub groups": {
			configPath: "testdata/TestDataUser/default.tf",
			init: func(t *testing.T, m *iam.Mock, testData testDataForUser) {
				expectGetUserMaxAuthGranSubGroups(t, m, testData, 3, maxSupportedGroupNesting)
			},
			mockData: basicUserTestDataMaxGroups,
		},
		"error - max amount of sub groups + 1": {
			configPath: "testdata/TestDataUser/default.tf",
			init: func(t *testing.T, m *iam.Mock, testData testDataForUser) {
				expectGetUserMaxAuthGranSubGroups(t, m, testData, 1, maxSupportedGroupNesting+1)
			},
			error:    regexp.MustCompile("unsupported subgroup depth"),
			mockData: basicUserTestDataMaxGroups,
		},
		"error - missing cidr_block_id": {
			configPath: "testdata/TestDataUser/missing_ui_identity.tf",
			error:      regexp.MustCompile("Missing required argument"),
			mockData:   basicUserTestData,
		},
		"error - GetUser call failed": {
			configPath: "testdata/TestDataUser/default.tf",
			init: func(t *testing.T, m *iam.Mock, user testDataForUser) {
				getUserReq := iam.GetUserRequest{IdentityID: user.uiIdentityID, Actions: true, AuthGrants: true, Notifications: true}
				m.On("GetUser", mock.Anything, getUserReq).Return(nil, errors.New("test error"))
			},
			error:    regexp.MustCompile("test error"),
			mockData: basicUserTestData,
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
							Check:       checkUserAttrs(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectGetUser(_ *testing.T, client *iam.Mock, data testDataForUser, times int) {
	getUserReq := iam.GetUserRequest{
		IdentityID:    data.uiIdentityID,
		Actions:       true,
		AuthGrants:    true,
		Notifications: true,
	}

	user := iam.User{
		IdentityID:                         data.uiIdentityID,
		IsLocked:                           data.isLocked,
		LastLoginDate:                      data.lastLoginDate,
		PasswordExpiryDate:                 data.passwordExpiryDate,
		TFAConfigured:                      data.tfaConfigured,
		EmailUpdatePending:                 data.emailUpdatePending,
		AccountID:                          data.accountID,
		AdditionalAuthenticationConfigured: data.additionalAuthenticationConfigured,
		UserBasicInfo: iam.UserBasicInfo{
			FirstName:                data.firstName,
			LastName:                 data.lastName,
			UserName:                 data.uiUserName,
			Email:                    data.email,
			Phone:                    data.phone,
			TimeZone:                 data.timeZone,
			JobTitle:                 data.jobTitle,
			TFAEnabled:               data.tfaEnabled,
			SecondaryEmail:           data.secondaryEmail,
			MobilePhone:              data.mobilePhone,
			Address:                  data.address,
			City:                     data.city,
			State:                    data.state,
			ZipCode:                  data.zipCode,
			Country:                  data.country,
			ContactType:              data.contactType,
			PreferredLanguage:        data.preferredLanguage,
			SessionTimeOut:           ptr.To(int(data.sessionTimeout)),
			AdditionalAuthentication: iam.Authentication(data.additionalAuthentication),
		},
		Actions: &iam.UserActions{
			APIClient:        data.actions.apiClient,
			Delete:           data.actions.delete,
			Edit:             data.actions.edit,
			IsCloneable:      data.actions.isCloneable,
			ResetPassword:    data.actions.resetPassword,
			ThirdPartyAccess: data.actions.thirdPartyAccess,
		},
		Notifications: iam.UserNotifications{
			Options: iam.UserNotificationOptions{
				NewUser:                   data.notifications.options.newUserNotification,
				PasswordExpiry:            data.notifications.options.passwordExpiry,
				Proactive:                 data.notifications.options.proactive,
				Upgrade:                   data.notifications.options.upgrade,
				APIClientCredentialExpiry: data.notifications.options.apiClientCredentialExpiryNotification,
			},
			EnableEmail: data.notifications.enableEmailNotifications,
		},
	}
	userAuthGrantList := make([]iam.AuthGrant, 0, len(data.authGrants))
	for _, authGrant := range data.authGrants {
		userAuthGrant := iam.AuthGrant{
			GroupID:         authGrant.groupID,
			GroupName:       authGrant.groupName,
			IsBlocked:       authGrant.isBlocked,
			RoleDescription: authGrant.roleDescription,
			RoleID:          ptr.To(int(authGrant.roleID)),
			RoleName:        authGrant.roleName,
		}
		userAuthGrantList = append(userAuthGrantList, userAuthGrant)
	}
	user.AuthGrants = userAuthGrantList

	client.On("GetUser", mock.Anything, getUserReq).Return(&user, nil).Times(times)
}
func expectGetUserMaxAuthGranSubGroups(_ *testing.T, client *iam.Mock, data testDataForUser, times, subGroupsDepth int) {
	getUserReq := iam.GetUserRequest{
		IdentityID:    data.uiIdentityID,
		Actions:       true,
		AuthGrants:    true,
		Notifications: true,
	}

	user := iam.User{
		IdentityID:                         data.uiIdentityID,
		IsLocked:                           data.isLocked,
		LastLoginDate:                      data.lastLoginDate,
		PasswordExpiryDate:                 data.passwordExpiryDate,
		TFAConfigured:                      data.tfaConfigured,
		EmailUpdatePending:                 data.emailUpdatePending,
		AccountID:                          data.accountID,
		AdditionalAuthenticationConfigured: data.additionalAuthenticationConfigured,
		UserBasicInfo: iam.UserBasicInfo{
			FirstName:                data.firstName,
			LastName:                 data.lastName,
			UserName:                 data.uiUserName,
			Email:                    data.email,
			Phone:                    data.phone,
			TimeZone:                 data.timeZone,
			JobTitle:                 data.jobTitle,
			TFAEnabled:               data.tfaEnabled,
			SecondaryEmail:           data.secondaryEmail,
			MobilePhone:              data.mobilePhone,
			Address:                  data.address,
			City:                     data.city,
			State:                    data.state,
			ZipCode:                  data.zipCode,
			Country:                  data.country,
			ContactType:              data.contactType,
			PreferredLanguage:        data.preferredLanguage,
			SessionTimeOut:           ptr.To(int(data.sessionTimeout)),
			AdditionalAuthentication: iam.Authentication(data.additionalAuthentication),
		},
		Actions: &iam.UserActions{
			APIClient:        data.actions.apiClient,
			Delete:           data.actions.delete,
			Edit:             data.actions.edit,
			IsCloneable:      data.actions.isCloneable,
			ResetPassword:    data.actions.resetPassword,
			ThirdPartyAccess: data.actions.thirdPartyAccess,
		},
		Notifications: iam.UserNotifications{
			Options: iam.UserNotificationOptions{
				NewUser:                   data.notifications.options.newUserNotification,
				PasswordExpiry:            data.notifications.options.passwordExpiry,
				Proactive:                 data.notifications.options.proactive,
				Upgrade:                   data.notifications.options.upgrade,
				APIClientCredentialExpiry: data.notifications.options.apiClientCredentialExpiryNotification,
			},
			EnableEmail: data.notifications.enableEmailNotifications,
		},
	}
	userAuthGrantList := make([]iam.AuthGrant, 0, len(data.authGrants))
	for _, authGrant := range data.authGrants {
		userAuthGrant := iam.AuthGrant{
			GroupID:         authGrant.groupID,
			GroupName:       authGrant.groupName,
			IsBlocked:       authGrant.isBlocked,
			RoleDescription: authGrant.roleDescription,
			RoleID:          ptr.To(int(authGrant.roleID)),
			RoleName:        authGrant.roleName,
			Subgroups:       generateMaxDepthSubGroupsAPIResponse(subGroupsDepth),
		}
		userAuthGrantList = append(userAuthGrantList, userAuthGrant)
	}
	user.AuthGrants = userAuthGrantList

	client.On("GetUser", mock.Anything, getUserReq).Return(&user, nil).Times(times)
}

func checkUserAttrs(data testDataForUser) resource.TestCheckFunc {
	name := "data.akamai_iam_user.test"
	checksFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "ui_identity_id", data.uiIdentityID),
	}
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "account_id", data.accountID))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "additional_authentication", data.additionalAuthentication))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "additional_authentication_configured", strconv.FormatBool(data.additionalAuthenticationConfigured)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "address", data.address))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "city", data.city))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "contact_type", data.contactType))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "state", data.state))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "country", data.country))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "email", data.email))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "email_update_pending", strconv.FormatBool(data.emailUpdatePending)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "first_name", data.firstName))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "is_locked", strconv.FormatBool(data.isLocked)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "job_title", data.jobTitle))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "last_login_date", date.FormatRFC3339Nano(data.lastLoginDate)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "last_name", data.lastName))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "mobile_phone", data.mobilePhone))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "password_expiry_date", date.FormatRFC3339Nano(data.passwordExpiryDate)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "phone", data.phone))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "preferred_language", data.preferredLanguage))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "secondary_email", data.secondaryEmail))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "session_timeout", strconv.FormatInt(data.sessionTimeout, 10)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "tfa_configured", strconv.FormatBool(data.tfaConfigured)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "tfa_enabled", strconv.FormatBool(data.tfaEnabled)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "time_zone", data.timeZone))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "ui_user_name", data.uiUserName))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "zip_code", data.zipCode))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "actions.delete", strconv.FormatBool(data.actions.delete)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "actions.api_client", strconv.FormatBool(data.actions.apiClient)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "actions.edit", strconv.FormatBool(data.actions.edit)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "actions.is_cloneable", strconv.FormatBool(data.actions.isCloneable)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "actions.reset_password", strconv.FormatBool(data.actions.resetPassword)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "actions.third_party_access", strconv.FormatBool(data.actions.thirdPartyAccess)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "auth_grants.#", strconv.Itoa(1)))
	for i, authGrant := range data.authGrants {
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("auth_grants.%d.group_id", i), strconv.FormatInt(authGrant.groupID, 10)))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("auth_grants.%d.group_name", i), authGrant.groupName))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("auth_grants.%d.is_blocked", i), strconv.FormatBool(authGrant.isBlocked)))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("auth_grants.%d.role_id", i), strconv.FormatInt(authGrant.roleID, 10)))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("auth_grants.%d.role_name", i), authGrant.roleName))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("auth_grants.%d.role_description", i), authGrant.roleDescription))
		if authGrant.subgroup != nil && len(authGrant.subgroup) > 0 {
			checksFuncs = append(checksFuncs, generateAggregateTestCheckFuncsForMaxAuthGrantSubGroups(i, 2, maxSupportedGroupNesting))
		}
	}

	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "notifications.enable_email_notifications", strconv.FormatBool(data.notifications.enableEmailNotifications)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "notifications.options.api_client_credential_expiry_notification", strconv.FormatBool(data.notifications.options.apiClientCredentialExpiryNotification)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "notifications.options.new_user_notification", strconv.FormatBool(data.notifications.options.newUserNotification)))
	checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, "notifications.options.password_expiry", strconv.FormatBool(data.notifications.options.passwordExpiry)))
	for i, proactiveElement := range data.notifications.options.proactive {
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("notifications.options.proactive.%d", i), proactiveElement))
	}
	for i, upgradeElement := range data.notifications.options.upgrade {
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("notifications.options.upgrade.%d", i), upgradeElement))
	}

	return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
}

func generateMaxDepthSubGroupsAuthGrantData(depth int) []*authGrantData {
	var subgroups []*authGrantData
	for groupID := depth; groupID > 1; groupID-- {

		authGrant := authGrantData{
			groupID:         int64(groupID),
			groupName:       fmt.Sprintf("group%d", groupID),
			isBlocked:       false,
			roleID:          int64(groupID),
			roleName:        fmt.Sprintf("role%d", groupID),
			roleDescription: fmt.Sprintf("roleDesc%d", groupID),
			subgroup:        subgroups,
		}
		subgroups = append([]*authGrantData{}, &authGrant)
	}
	return subgroups
}

func generateMaxDepthSubGroupsAPIResponse(depth int) []iam.AuthGrant {
	var subgroups []iam.AuthGrant
	for groupID := depth; groupID > 1; groupID-- {

		authGrant := iam.AuthGrant{
			GroupID:         int64(groupID),
			GroupName:       fmt.Sprintf("group%d", groupID),
			IsBlocked:       false,
			RoleID:          ptr.To(groupID),
			RoleName:        fmt.Sprintf("role%d", groupID),
			RoleDescription: fmt.Sprintf("roleDesc%d", groupID),
			Subgroups:       subgroups,
		}
		subgroups = []iam.AuthGrant{authGrant}
	}
	return subgroups
}

func generateAggregateTestCheckFuncsForMaxAuthGrantSubGroups(authGrantElement, min, max int) resource.TestCheckFunc {
	var testCases []resource.TestCheckFunc
	path := "sub_groups.0"
	for i := min; i < max; i++ {
		testCases = append(testCases, resource.TestCheckResourceAttr("data.akamai_iam_user.test", fmt.Sprintf("auth_grants.%d.%s.group_id", authGrantElement, path), strconv.Itoa(i)))
		testCases = append(testCases, resource.TestCheckResourceAttr("data.akamai_iam_user.test", fmt.Sprintf("auth_grants.%d.%s.group_name", authGrantElement, path), fmt.Sprintf("group%d", i)))
		testCases = append(testCases, resource.TestCheckResourceAttr("data.akamai_iam_user.test", fmt.Sprintf("auth_grants.%d.%s.role_id", authGrantElement, path), strconv.Itoa(i)))
		testCases = append(testCases, resource.TestCheckResourceAttr("data.akamai_iam_user.test", fmt.Sprintf("auth_grants.%d.%s.role_name", authGrantElement, path), fmt.Sprintf("role%d", i)))
		testCases = append(testCases, resource.TestCheckResourceAttr("data.akamai_iam_user.test", fmt.Sprintf("auth_grants.%d.%s.role_description", authGrantElement, path), fmt.Sprintf("roleDesc%d", i)))
		testCases = append(testCases, resource.TestCheckResourceAttr("data.akamai_iam_user.test", fmt.Sprintf("auth_grants.%d.%s.is_blocked", authGrantElement, path), "false"))
		path = path + ".sub_groups.0"
	}

	return resource.ComposeAggregateTestCheckFunc(testCases...)
}
