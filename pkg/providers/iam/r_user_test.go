package iam

import (
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResUserLifecycle(t *testing.T) {
	t.Parallel()

	// A lifecycle test creates one user resource and takes it through 7 updates before destroying it. These tests only
	// vary in the first step which defines the circumstances of creation. The validations and expected service API
	// behavior for each test are otherwise identical.
	//
	// These specific variations in creation and the update steps were chosen ensure all attributes can be changed;
	// optional attributes can be added, changed, and removed; and to ensure correct interpretations of defaults.

	// encoded slice of one zero-value AuthGrant
	OneAuthGrantJSONA := `[{"groupId":0,"groupName":"A","isBlocked":false,"roleDescription":"","roleName":""}]`
	OneAuthGrantJSONB := `[{"groupId":0,"groupName":"B","isBlocked":false,"roleDescription":"","roleName":""}]`

	type TestState struct {
		Provider   *provider
		Client     *IAM
		User       iam.User
		UserExists bool
	}

	// Any function that configures the mock IAM service behavior
	type TestSetupFunc = func(*TestState)

	// Describes a Lifecycle test's expectations for the creation step and transition to the first update
	type TestVariant struct {
		// Determines the test name and which fixture to load for the creation step
		Name string

		// Validates the resource state after creation
		Check resource.TestCheckFunc

		// Configures the mock before creation step
		Setup TestSetupFunc

		// Configures the mock before the first update step
		Transition TestSetupFunc
	}

	// Setup the standard GetUser expectation
	ExpectGetUser := func(State *TestState) {
		req := iam.GetUserRequest{
			IdentityID:    "test uiIdentityId",
			AuthGrants:    true,
			Notifications: true,
		}

		call := State.Client.On("GetUser", mock.Anything, req)
		call.Run(func(mock.Arguments) {
			if !State.UserExists {
				call.Return(nil, fmt.Errorf("Not Found"))
			}

			u := CopyUser(State.User)
			call.Return(&u, nil)
		})
	}

	ExpectCreateUser := func(State *TestState, User iam.User) {
		req := iam.CreateUserRequest{
			User:      CopyBasicUser(User.UserBasicInfo),
			SendEmail: true,
		}
		req.Notifications = User.Notifications
		req.AuthGrants = User.AuthGrants

		call := State.Client.On("CreateUser", mock.Anything, req).Once()
		call.Run(func(mock.Arguments) {
			res := CopyUser(User)
			State.User = CopyUser(User)
			State.UserExists = true
			call.Return(&res, nil)
		})
	}

	ExpectUpdateUserInfo := func(State *TestState, User iam.UserBasicInfo) {
		req := iam.UpdateUserInfoRequest{
			IdentityID: "test uiIdentityId",
			User:       User,
		}

		call := State.Client.On("UpdateUserInfo", mock.Anything, req).Once()
		call.Run(func(mock.Arguments) {
			res := CopyBasicUser(User)
			State.User.UserBasicInfo = CopyBasicUser(User)
			call.Return(&res, nil)
		})
	}

	ExpectUpdateUserNotifications := func(State *TestState, Notifications iam.UserNotifications) {
		req := iam.UpdateUserNotificationsRequest{
			IdentityID:    "test uiIdentityId",
			Notifications: Notifications,
		}

		call := State.Client.On("UpdateUserNotifications", mock.Anything, req).Once()
		call.Run(func(mock.Arguments) {
			n := CopyUserNotifications(Notifications)
			res := CopyUserNotifications(Notifications)
			State.User.Notifications = &n
			call.Return(&res, nil)
		})
	}

	ExpectUpdateAuthGrants := func(State *TestState, AuthGrants []iam.AuthGrant) {
		req := iam.UpdateUserAuthGrantsRequest{
			IdentityID: "test uiIdentityId",
			AuthGrants: AuthGrants,
		}

		call := State.Client.On("UpdateUserAuthGrants", mock.Anything, req).Once()
		call.Run(func(mock.Arguments) {
			State.User.AuthGrants = CopyAuthGrants(AuthGrants)
			call.Return(CopyAuthGrants(AuthGrants), nil)
		})
	}

	ExpectRemoveUser := func(State *TestState) {
		req := iam.RemoveUserRequest{IdentityID: "test uiIdentityId"}

		call := State.Client.On("RemoveUser", mock.Anything, req).Once()
		call.Run(func(mock.Arguments) {
			State.User = iam.User{}
			State.UserExists = false
			call.Return(nil)
		})
	}

	authGrants := func(name string) []iam.AuthGrant {
		return []iam.AuthGrant{{GroupName: name}}
	}

	mkUser := func(Basic iam.UserBasicInfo, Notifications *iam.UserNotifications, AuthGrants []iam.AuthGrant) iam.User {
		User := iam.User{
			IdentityID:         "test uiIdentityId",
			LastLoginDate:      "last login",
			PasswordExpiryDate: "password expired after",
			IsLocked:           true,
			TFAConfigured:      true,
			EmailUpdatePending: true,
		}
		User.UserBasicInfo = Basic
		User.Notifications = Notifications
		User.AuthGrants = AuthGrants
		return User
	}

	// Notifications variation A
	notifA := func() *iam.UserNotifications {
		return &iam.UserNotifications{EnableEmail: true}
	}

	// Notifications variation C
	notifC := func() *iam.UserNotifications {
		return &iam.UserNotifications{
			EnableEmail: true,
			Options: &iam.UserNotificationOptions{
				NewUser:        true,
				PasswordExpiry: true,
				Proactive:      []string{"issues product"},
				Upgrade:        []string{"upgrades product"},
			},
		}
	}

	// minimum user attributes variation A
	minUserA := func() iam.UserBasicInfo {
		return iam.UserBasicInfo{
			FirstName:  "first name A",
			LastName:   "last name A",
			Email:      "email-A@akamai.net",
			Phone:      "phone A",
			TFAEnabled: true,
			Country:    "country A",
		}
	}

	// All basic user info variation A
	allUserA := func() iam.UserBasicInfo {
		SessionTimeOut := 1
		return iam.UserBasicInfo{
			FirstName:         "first name A",
			LastName:          "last name A",
			UserName:          "user name A",
			Email:             "email-A@akamai.net",
			Phone:             "phone A",
			TimeZone:          "Timezone A",
			JobTitle:          "job title A",
			TFAEnabled:        true,
			SecondaryEmail:    "secondary-email-A@akamai.net",
			MobilePhone:       "mobile phone A",
			Address:           "123 A Street",
			City:              "A-Town",
			State:             "state A",
			ZipCode:           "zip A",
			Country:           "country A",
			ContactType:       "contact type A",
			PreferredLanguage: "language A",
			SessionTimeOut:    &SessionTimeOut,
		}
	}

	// minimum user attributes variation B
	minUserB := func() iam.UserBasicInfo {
		return iam.UserBasicInfo{
			FirstName: "first name B",
			LastName:  "last name B",
			Email:     "email-B@akamai.net",
			Phone:     "phone B",
			Country:   "country B",
		}
	}

	// All basic user info variation B
	allUserB := func() iam.UserBasicInfo {
		SessionTimeout := 2
		return iam.UserBasicInfo{
			FirstName:         "first name B",
			LastName:          "last name B",
			UserName:          "user name B",
			Email:             "email-B@akamai.net",
			Phone:             "phone B",
			TimeZone:          "Timezone B",
			JobTitle:          "job title B",
			TFAEnabled:        false,
			SecondaryEmail:    "secondary-email-B@akamai.net",
			MobilePhone:       "mobile phone B",
			Address:           "123 B Street",
			City:              "B-Town",
			State:             "state B",
			ZipCode:           "zip B",
			Country:           "country B",
			ContactType:       "contact type B",
			PreferredLanguage: "language B",
			SessionTimeOut:    &SessionTimeout,
		}
	}

	// Compose a resource.TestCheckFunc that verifies all attributes of the given user
	CheckState := func(User iam.User) resource.TestCheckFunc {
		if User.SessionTimeOut == nil {
			SessionTimeout := 0
			User.SessionTimeOut = &SessionTimeout
		}

		var AuthGrantsJSON string
		var ProductChecks []resource.TestCheckFunc
		if len(User.AuthGrants) > 0 {
			switch User.AuthGrants[0].GroupName {
			case "A":
				AuthGrantsJSON = OneAuthGrantJSONA
			case "B":
				AuthGrantsJSON = OneAuthGrantJSONB
			default:
				panic("unknown auth grant group name")
			}
		}

		EnableNotifications := "false"
		PasswordExpiry := "false"
		NewUsers := "false"
		if User.Notifications != nil {
			EnableNotifications = fmt.Sprintf("%t", User.Notifications.EnableEmail)

			if User.Notifications.Options != nil {
				NewUsers = fmt.Sprintf("%t", User.Notifications.Options.NewUser)
				PasswordExpiry = fmt.Sprintf("%t", User.Notifications.Options.PasswordExpiry)

				for _, p := range User.Notifications.Options.Proactive {
					chk := resource.TestCheckTypeSetElemAttr("akamai_iam_user.test", "subscribe_product_issues.*", p)
					ProductChecks = append(ProductChecks, chk)
				}

				for _, p := range User.Notifications.Options.Upgrade {
					chk := resource.TestCheckTypeSetElemAttr("akamai_iam_user.test", "subscribe_product_upgrades.*", p)
					ProductChecks = append(ProductChecks, chk)
				}
			}
		}

		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("akamai_iam_user.test", "id", "test uiIdentityId"),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "first_name", User.FirstName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "last_name", User.LastName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "email", User.Email),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "country", User.Country),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "phone", User.Phone),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "enable_tfa", fmt.Sprintf("%t", User.TFAEnabled)),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "contact_type", User.ContactType),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "user_name", User.UserName),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "job_title", User.JobTitle),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "time_zone", User.TimeZone),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "secondary_email", User.SecondaryEmail),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "mobile_phone", User.MobilePhone),
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
			resource.TestCheckResourceAttr("akamai_iam_user.test", "auth_grants_json", AuthGrantsJSON),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "enable_notifications", EnableNotifications),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "subscribe_new_users", NewUsers),
			resource.TestCheckResourceAttr("akamai_iam_user.test", "subscribe_password_expiration", PasswordExpiry),
		}

		return resource.ComposeAggregateTestCheckFunc(append(checks, ProductChecks...)...)
	}

	// Setup each step by Asserting mock expectations then swap in a new mock
	InitStep := func(t *testing.T, State *TestState) {
		if State.Client != nil {
			if !State.Client.AssertExpectations(t) {
				t.FailNow()
			}
		}

		State.Client = &IAM{}
		State.Client.Test(test.TattleT{T: t})
		ExpectGetUser(State)
		State.Provider.SetIAM(State.Client)
	}

	// Drive a Lifecycle test case
	AssertLifecycle := func(t *testing.T, tv TestVariant) {
		t.Helper()
		fixturePrefix := fmt.Sprintf("testdata/%s", t.Name())

		t.Run(tv.Name, func(t *testing.T) {
			t.Helper()
			t.Parallel()

			State := TestState{Provider: &provider{}}
			State.Provider.SetCache(metaCache{})

			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: State.Provider.ProviderFactories(),
				Steps: []resource.TestStep{
					{ // Step 0 - Varies
						Config: test.Fixture("%s/step00-%s.tf", fixturePrefix, tv.Name),
						PreConfig: func() {
							InitStep(t, &State)
							tv.Setup(&State)
						},
						Check: tv.Check,
					},
					{ // Step 1 - Minimum user attributes variation B
						Config: test.Fixture("%s/step01.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							tv.Transition(&State)
						},
						Check: CheckState(mkUser(minUserB(), nil, authGrants("B"))),
					},
					{ // Step 2 - All user attributes variation B
						Config: test.Fixture("%s/step02.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							ExpectUpdateUserInfo(&State, allUserB())
						},
						Check: CheckState(mkUser(allUserB(), nil, authGrants("B"))),
					},
					{ // Step 3 - All user attributes variation B, notifications C, grants
						Config: test.Fixture("%s/step03.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							ExpectUpdateUserNotifications(&State, *notifC())
							// ExpectUpdateAuthGrants(&State, authGrants("A"))
						},
						Check: CheckState(mkUser(allUserB(), notifC(), authGrants("B"))),
					},
					{ // Step 4 - All user attributes variation B, notifications B
						Config: test.Fixture("%s/step04.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							ExpectUpdateUserNotifications(&State, iam.UserNotifications{})
							// ExpectUpdateAuthGrants(&State, authGrants("A"))
						},
						Check: CheckState(mkUser(allUserB(), nil, authGrants("B"))),
					},
					{ // Step 5 - All user attributes variation B, notifications A
						Config: test.Fixture("%s/step05.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							ExpectUpdateUserNotifications(&State, iam.UserNotifications{EnableEmail: true})
						},
						Check: CheckState(mkUser(allUserB(), notifA(), authGrants("B"))),
					},
					{ // Step 6 - All user attributes variation B, notifications B
						Config: test.Fixture("%s/step06.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							ExpectUpdateUserNotifications(&State, iam.UserNotifications{})
						},
						Check: CheckState(mkUser(allUserB(), nil, authGrants("B"))),
					},
					{ // Step 7 - minimum user attributes variation B
						Config: test.Fixture("%s/step07.tf", fixturePrefix),
						PreConfig: func() {
							InitStep(t, &State)
							ExpectUpdateUserInfo(&State, minUserB())
							ExpectRemoveUser(&State)
						},
						Check: CheckState(mkUser(minUserB(), nil, authGrants("B"))),
					},
				}, // Steps
			}) // resource.UnitTest()

			State.Client.AssertExpectations(t)
		}) // t.Run()
	} // AssertLifecycle

	AssertLifecycle(t, TestVariant{
		Name:  "minimum basic info A",
		Check: CheckState(mkUser(minUserA(), nil, authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(minUserA(), nil, authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})

	AssertLifecycle(t, TestVariant{
		Name:  "all basic info A",
		Check: CheckState(mkUser(allUserA(), nil, authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(allUserA(), nil, authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})

	AssertLifecycle(t, TestVariant{
		Name:  "all basic info A with notifications C",
		Check: CheckState(mkUser(allUserA(), notifC(), authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(allUserA(), notifC(), authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateUserNotifications(State, iam.UserNotifications{})
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})

	AssertLifecycle(t, TestVariant{
		Name:  "all basic info A with grants",
		Check: CheckState(mkUser(allUserA(), nil, authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(allUserA(), nil, authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})

	AssertLifecycle(t, TestVariant{
		Name:  "all basic info A with notifications C and grants",
		Check: CheckState(mkUser(allUserA(), notifC(), authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(allUserA(), notifC(), authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateUserNotifications(State, iam.UserNotifications{})
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})

	AssertLifecycle(t, TestVariant{
		Name:  "all basic info A with notifications A",
		Check: CheckState(mkUser(allUserA(), notifA(), authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(allUserA(), notifA(), authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateUserNotifications(State, iam.UserNotifications{})
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})

	AssertLifecycle(t, TestVariant{
		Name:  "all basic info A with notifications B",
		Check: CheckState(mkUser(allUserA(), nil, authGrants("A"))),
		Setup: func(State *TestState) {
			ExpectCreateUser(State, mkUser(allUserA(), nil, authGrants("A")))
		},
		Transition: func(State *TestState) {
			ExpectUpdateUserInfo(State, minUserB())
			ExpectUpdateAuthGrants(State, authGrants("B"))
		},
	})
}

func CopyBasicUser(User iam.UserBasicInfo) iam.UserBasicInfo {
	if User.SessionTimeOut != nil {
		SessionTimeOut := *User.SessionTimeOut
		User.SessionTimeOut = &SessionTimeOut
	}

	return User
}

// Deeply duplicate the given iam.User
func CopyUser(User iam.User) iam.User {
	if User.SessionTimeOut != nil {
		SessionTimeOut := *User.SessionTimeOut
		User.SessionTimeOut = &SessionTimeOut
	}

	if User.AuthGrants != nil {
		AuthGrants := []iam.AuthGrant{}

		for _, ag := range User.AuthGrants {
			AuthGrants = append(AuthGrants, CopyAuthGrant(ag))
		}

		User.AuthGrants = AuthGrants
	}

	if User.Notifications != nil {
		Notifications := CopyUserNotifications(*User.Notifications)
		User.Notifications = &Notifications
	}

	return User
}

// Deeply duplicate a UserNotifications
func CopyUserNotifications(Notifications iam.UserNotifications) iam.UserNotifications {
	if Notifications.Options != nil {
		Options := *Notifications.Options

		if Options.Proactive != nil {
			Proactive := make([]string, len(Options.Proactive))
			copy(Proactive, Options.Proactive)
			Options.Proactive = Proactive
		}

		if Options.Upgrade != nil {
			Upgrade := make([]string, len(Options.Upgrade))
			copy(Upgrade, Options.Upgrade)
			Options.Upgrade = Upgrade
		}

		Notifications.Options = &Options
	}

	return Notifications
}

func CopyAuthGrants(AuthGrants []iam.AuthGrant) []iam.AuthGrant {
	var cp []iam.AuthGrant

	for _, ag := range AuthGrants {
		cp = append(cp, CopyAuthGrant(ag))
	}

	return cp
}

// Deeply duplicate the given iam.AuthGrant
func CopyAuthGrant(AuthGrant iam.AuthGrant) iam.AuthGrant {
	if AuthGrant.RoleID != nil {
		RoleID := *AuthGrant.RoleID
		AuthGrant.RoleID = &RoleID
	}

	if AuthGrant.Subgroups != nil {
		Subgroups := []iam.AuthGrant{}
		for _, ag := range AuthGrant.Subgroups {
			Subgroups = append(Subgroups, CopyAuthGrant(ag))
		}
	}

	return AuthGrant
}
