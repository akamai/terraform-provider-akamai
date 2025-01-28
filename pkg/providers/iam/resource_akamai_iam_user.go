package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var northAmerica = regexp.MustCompile(`^(\s*\+\s*1\s*|\s*1\s*|\s*)(|-|\.|\\|\()([1-9][0-9]{2})(|-|\.|\\|\))\s*([0-9]{3})(|-|\.|\\|\s)([0-9]{4})(|-|\.|\\|\sx?)([0-9]{1,30})?$`)
var international = regexp.MustCompile(`^\+[02-9][\d\s\-]{0,40}$`)

func resourceIAMUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a user in your account.",
		CreateContext: resourceIAMUserCreate,
		ReadContext:   resourceIAMUserRead,
		UpdateContext: resourceIAMUserUpdate,
		DeleteContext: resourceIAMUserDelete,
		CustomizeDiff: customdiff.All(customizePasswordDiff, customizeNotificationDiff),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Inputs - Required
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's first name.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's surname.",
			},
			"email": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The user's email address.",
				StateFunc:        stateEmail,
				DiffSuppressFunc: suppressEmail,
			},
			"country": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "As part of the user's location, the value can be any that are available from the view-supported-countries operation.",
			},
			"phone": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The user's main phone number.",
				DiffSuppressFunc: suppressPhone,
				StateFunc:        statePhone,
				ValidateFunc:     validatePhone,
			},
			"enable_tfa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether two-factor authentication is allowed.",
			},
			"enable_mfa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether multi-factor authentication is allowed.",
			},
			"auth_grants_json": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "A user's per-group role assignments, in JSON form.",
				ValidateDiagFunc: validateAuthGrantsJS,
				DiffSuppressFunc: suppressAuthGrantsJS,
				StateFunc:        stateAuthGrantsJS,
			},

			// Inputs - Optional
			"contact_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To help characterize the user, the value can be any that are available from the view-contact-types operation.",
			},
			"job_title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's position at your company.",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user's time zone. The value can be any that are available from the view-time-zones operation.",
			},
			"secondary_email": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The user's secondary email address.",
				StateFunc:        stateEmail,
				DiffSuppressFunc: suppressEmail,
			},
			"mobile_phone": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The user's mobile phone number.",
				DiffSuppressFunc: suppressPhone,
				StateFunc:        statePhone,
				ValidateFunc:     validatePhone,
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user's street address.",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's city.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's state.",
			},
			"zip_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's five-digit ZIP code.",
			},
			"preferred_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The value can be any that are available from the view-languages operation.",
			},
			"session_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of seconds it takes for the user's Control Center session to time out if there hasn't been any activity.",
			},

			// Purely computed
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A user's `loginId`. Typically, a user's email address.",
			},
			"last_login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the user last logged in.",
			},
			"password_expired_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a user's password expires.",
			},
			"tfa_configured": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether two-factor authentication is configured.",
			},
			"email_update_pending": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether email update is pending.",
			},
			"lock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to block a user account.",
				Default:     false,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "New password for a user.",
				Optional:    true,
				Sensitive:   true,
			},
			"user_notifications": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies email notifications the user receives for products.",
				Computed:    true,
				MaxItems:    1, // Ensure only one notification configuration can be set
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_client_credential_expiry_notification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enables notifications for expiring API client credentials.",
						},
						"new_user_notification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enables notifications for group administrators when the user creates other new users.",
						},
						"password_expiry": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enables notifications for expiring passwords.",
						},
						"proactive": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Products for which the user gets notifications for service issues.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"upgrade": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Products for which the user receives notifications for upgrades.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enable_email_notifications": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enables email notifications.",
						},
					},
				},
			},
		},
	}
}

func resourceIAMUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "resourceIAMUserCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Creating User")

	authGrantsJSON := []byte(d.Get("auth_grants_json").(string))

	var authGrants []iam.AuthGrantRequest
	if len(authGrantsJSON) > 0 {
		if err := json.Unmarshal(authGrantsJSON, &authGrants); err != nil {
			logger.Error("auth_grants is not valid", "error", err)
			return diag.Errorf("auth_grants is not valid: %s", err)
		}
	}

	authMethod, err := getAuthenticationMethod(d)
	if err != nil {
		return diag.FromErr(err)
	}

	basicUser := iam.UserBasicInfo{
		FirstName:                d.Get("first_name").(string),
		LastName:                 d.Get("last_name").(string),
		UserName:                 d.Get("user_name").(string),
		Email:                    d.Get("email").(string),
		Phone:                    d.Get("phone").(string),
		TimeZone:                 d.Get("time_zone").(string),
		JobTitle:                 d.Get("job_title").(string),
		TFAEnabled:               d.Get("enable_tfa").(bool),
		SecondaryEmail:           d.Get("secondary_email").(string),
		MobilePhone:              d.Get("mobile_phone").(string),
		Address:                  d.Get("address").(string),
		City:                     d.Get("city").(string),
		State:                    d.Get("state").(string),
		ZipCode:                  d.Get("zip_code").(string),
		Country:                  d.Get("country").(string),
		ContactType:              d.Get("contact_type").(string),
		PreferredLanguage:        d.Get("preferred_language").(string),
		AdditionalAuthentication: iam.Authentication(authMethod),
	}

	if st, ok := d.GetOk("session_timeout"); ok {
		sessionTimeout := st.(int)
		basicUser.SessionTimeOut = &sessionTimeout
	}

	userRequest := iam.CreateUserRequest{
		UserBasicInfo: basicUser,
		AuthGrants:    authGrants,
		SendEmail:     true,
	}

	// Get user notifications if provided
	userNotifications, err := getUserNotifications(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if userNotifications != nil {
		userRequest.Notifications = userNotifications
	}

	user, err := client.CreateUser(ctx, userRequest)
	if err != nil {
		logger.Error("failed to create user", "error", err)
		return diag.Errorf("failed to create user: %s\n%s", err, resourceIAMUserErrorAdvice(err))
	}

	err = manageUserPassword(ctx, d, client, user.IdentityID)
	if err != nil {
		logger.Errorf("failed to set user password", "error", err)
		return diag.Errorf("failed to set user password: %s", err)
	}

	// lock the user's account
	lock, err := tf.GetBoolValue("lock", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if lock {
		if err = client.LockUser(ctx, iam.LockUserRequest{IdentityID: user.IdentityID}); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(user.IdentityID)
	return resourceIAMUserRead(ctx, d, m)
}

func manageUserPassword(ctx context.Context, d *schema.ResourceData, client iam.IAM, ID string) error {
	password, err := tf.GetStringValue("password", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if password != "" {
		err = client.SetUserPassword(ctx, iam.SetUserPasswordRequest{
			IdentityID:  ID,
			NewPassword: password,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func extractUserNotificationsData(notificationsData interface{}) (*iam.UserNotifications, error) {

	notificationsList, ok := notificationsData.([]interface{})
	if !ok {
		return nil, errors.New("user notifications data is not a valid list")
	}

	itemMap, ok := notificationsList[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("user notifications data item is not a valid map")
	}

	return &iam.UserNotifications{
		EnableEmail: itemMap["enable_email_notifications"].(bool),
		Options: iam.UserNotificationOptions{
			APIClientCredentialExpiry: itemMap["api_client_credential_expiry_notification"].(bool),
			NewUser:                   itemMap["new_user_notification"].(bool),
			PasswordExpiry:            itemMap["password_expiry"].(bool),
			Proactive:                 tf.InterfaceSliceToStringSlice(itemMap["proactive"].([]interface{})),
			Upgrade:                   tf.InterfaceSliceToStringSlice(itemMap["upgrade"].([]interface{})),
		},
	}, nil
}

func getUserNotifications(d *schema.ResourceData) (*iam.UserNotifications, error) {
	notificationsData, err := tf.GetListValue("user_notifications", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	return extractUserNotificationsData(notificationsData)
}

func getAuthenticationMethod(d *schema.ResourceData) (string, error) {
	enableTFA := d.Get("enable_tfa").(bool)
	enableMFA := d.Get("enable_mfa").(bool)

	if enableTFA && enableMFA {
		return "", errors.New("only one of 'enable_tfa' or 'enable_mfa' can be set")
	}

	if enableTFA {
		return "TFA", nil
	} else if enableMFA {
		return "MFA", nil
	}
	return "NONE", nil
}

func resourceIAMUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "resourceIAMUserRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading User")

	req := iam.GetUserRequest{
		IdentityID:    d.Id(),
		AuthGrants:    true,
		Notifications: true,
	}

	user, err := client.GetUser(ctx, req)
	if err != nil {
		logger.Error("failed to fetch user", "error", err)
		return diag.Errorf("failed to fetch user: %s", err)
	}

	if user.SessionTimeOut == nil {
		sessionTimeOut := 0
		user.SessionTimeOut = &sessionTimeOut
	}

	var authGrantsJSON []byte
	if len(user.AuthGrants) > 0 {
		authGrantsJSON, err = json.Marshal(user.AuthGrants)
		if err != nil {
			logger.Error("could not marshal AuthGrants", "error", err)
			return diag.Errorf("could not marshal AuthGrants: %s", err)
		}
	}

	userNotifications := []interface{}{map[string]interface{}{
		"enable_email_notifications":                user.Notifications.EnableEmail,
		"api_client_credential_expiry_notification": user.Notifications.Options.APIClientCredentialExpiry,
		"new_user_notification":                     user.Notifications.Options.NewUser,
		"password_expiry":                           user.Notifications.Options.PasswordExpiry,
		"proactive":                                 user.Notifications.Options.Proactive,
		"upgrade":                                   user.Notifications.Options.Upgrade,
	},
	}

	enableMFA := user.AdditionalAuthentication == "MFA"

	err = tf.SetAttrs(d, map[string]interface{}{
		"first_name":             user.FirstName,
		"last_name":              user.LastName,
		"user_name":              user.UserName,
		"email":                  user.Email,
		"phone":                  user.Phone,
		"time_zone":              user.TimeZone,
		"job_title":              user.JobTitle,
		"enable_tfa":             user.TFAEnabled,
		"enable_mfa":             enableMFA,
		"secondary_email":        user.SecondaryEmail,
		"mobile_phone":           user.MobilePhone,
		"address":                user.Address,
		"city":                   user.City,
		"state":                  user.State,
		"zip_code":               user.ZipCode,
		"country":                user.Country,
		"contact_type":           user.ContactType,
		"preferred_language":     user.PreferredLanguage,
		"last_login":             date.FormatRFC3339Nano(user.LastLoginDate),
		"password_expired_after": date.FormatRFC3339Nano(user.PasswordExpiryDate),
		"tfa_configured":         user.TFAConfigured,
		"email_update_pending":   user.EmailUpdatePending,
		"session_timeout":        *user.SessionTimeOut,
		"auth_grants_json":       stateAuthGrantsJS(string(authGrantsJSON)),
		"lock":                   user.IsLocked,
		"user_notifications":     userNotifications,
	})
	if err != nil {
		logger.Error("could not save attributes to state", "error", err)
		return diag.Errorf("could not save attributes to state: %s", err)
	}

	return nil
}

func resourceIAMUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "resourceIAMUserUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Updating User")

	if d.HasChange("email") {
		d.Partial(true)
		err := fmt.Errorf("cannot change email address")
		logger.Error("failed to update user", "error", err)
		return diag.Errorf("failed to update user: %s", err)
	}

	var needRead bool

	// Basic Info
	updateBasicInfo := d.HasChanges(
		"first_name",
		"last_name",
		"phone",
		"time_zone",
		"job_title",
		"enable_tfa",
		"enable_mfa",
		"secondary_email",
		"mobile_phone",
		"address",
		"city",
		"state",
		"zip_code",
		"country",
		"contact_type",
		"preferred_language",
		"session_timeout",
	)

	authMethod, err := getAuthenticationMethod(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if updateBasicInfo {
		basicUser := iam.UserBasicInfo{
			FirstName:                d.Get("first_name").(string),
			LastName:                 d.Get("last_name").(string),
			UserName:                 d.Get("user_name").(string),
			Email:                    d.Get("email").(string),
			Phone:                    d.Get("phone").(string),
			TimeZone:                 d.Get("time_zone").(string),
			JobTitle:                 d.Get("job_title").(string),
			TFAEnabled:               d.Get("enable_tfa").(bool),
			SecondaryEmail:           d.Get("secondary_email").(string),
			MobilePhone:              d.Get("mobile_phone").(string),
			Address:                  d.Get("address").(string),
			City:                     d.Get("city").(string),
			State:                    d.Get("state").(string),
			ZipCode:                  d.Get("zip_code").(string),
			Country:                  d.Get("country").(string),
			ContactType:              d.Get("contact_type").(string),
			PreferredLanguage:        d.Get("preferred_language").(string),
			AdditionalAuthentication: iam.Authentication(authMethod),
		}

		if st, ok := d.GetOk("session_timeout"); ok {
			sessionTimeout := st.(int)
			basicUser.SessionTimeOut = &sessionTimeout
		}

		req := iam.UpdateUserInfoRequest{
			IdentityID: d.Id(),
			User:       basicUser,
		}
		if _, err := client.UpdateUserInfo(ctx, req); err != nil {
			d.Partial(true)
			logger.Error("failed to update user", "error", err)
			return diag.Errorf("failed to update user: %s\n%s", err, resourceIAMUserErrorAdvice(err))
		}

		needRead = true
	}

	// AuthGrants
	if d.HasChange("auth_grants_json") {
		var authGrants []iam.AuthGrantRequest

		authGrantsJSON := []byte(d.Get("auth_grants_json").(string))
		if len(authGrantsJSON) > 0 {
			if err := json.Unmarshal(authGrantsJSON, &authGrants); err != nil {
				d.Partial(true)
				logger.Error("auth_grants is not valid", "error", err)
				return diag.Errorf("auth_grants is not valid: %s", err)
			}
		}

		req := iam.UpdateUserAuthGrantsRequest{
			IdentityID: d.Id(),
			AuthGrants: authGrants,
		}
		if _, err := client.UpdateUserAuthGrants(ctx, req); err != nil {
			d.Partial(true)
			logger.Error("failed to update user AuthGrants", "error", err)
			return diag.Errorf("failed to update user AuthGrants: %s", err)
		}

		needRead = true
	}

	// user notifications
	if d.HasChange("user_notifications") {
		notificationsData, err := tf.GetListValue("user_notifications", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}

		userNotifications, err := extractUserNotificationsData(notificationsData)
		if err != nil {
			return diag.Errorf("failed to extract user notifications: %s", err)
		}

		req := iam.UpdateUserNotificationsRequest{
			IdentityID:    d.Id(),
			Notifications: userNotifications,
		}

		if _, err := client.UpdateUserNotifications(ctx, req); err != nil {
			d.Partial(true)
			logger.Error("failed to update user notifications", "error", err)
			return diag.Errorf("failed to update user notifications: %s", err)
		}
		needRead = true
	}

	// password
	if d.HasChange("password") {
		err = manageUserPassword(ctx, d, client, d.Id())
		if err != nil {
			logger.Error("failed to set user password", "error", err)
			return diag.Errorf("failed to set user password: %s", err)
		}
		needRead = true
	}

	// lock the user
	if d.HasChange("lock") {
		lock, err := tf.GetBoolValue("lock", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}

		if lock {
			err = client.LockUser(ctx, iam.LockUserRequest{IdentityID: d.Id()})
		} else {
			err = client.UnlockUser(ctx, iam.UnlockUserRequest{IdentityID: d.Id()})
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if needRead {
		d.Partial(false)
		return resourceIAMUserRead(ctx, d, m)
	}

	d.Partial(false)
	return nil
}

func resourceIAMUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "resourceIAMUserDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Deleting User")

	if err := client.RemoveUser(ctx, iam.RemoveUserRequest{IdentityID: d.Id()}); err != nil {
		logger.Error("could not remove user", "error", err)
		return diag.Errorf("could not remove user: %s", err)
	}

	return nil
}

func resourceIAMUserErrorAdvice(e error) string {
	switch {
	case regexp.MustCompile(`\b(preferredLanguage|[pP]referred [lL]anguage)\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_supported_langs" data source to get possible values for "preferred_language"`

	case regexp.MustCompile(`\b(contactType|[cC]ontact [tT]ype)\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_contact_types" data source to get possible values for "contact_type"`

	case regexp.MustCompile(`\b[cC]ountry\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_countries" data source to get possible values for "country"`

	case regexp.MustCompile(`\b(sessionTimeOut|[sS]ession [tT]ime ?[oO]ut)\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_timeout_policies" data source to get possible values for "session_timeout"`

	case regexp.MustCompile(`\b[sS]tate\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_states" data source to get possible values for "state"`
	}

	return ""
}

func canonicalPhone(in string) string {
	if northAmerica.MatchString(in) {
		ph := northAmerica.FindStringSubmatch(in)
		if ph[9] == "" { // without extension
			return fmt.Sprintf("(%s) %s-%s", ph[3], ph[5], ph[7])
		}
		return fmt.Sprintf("(%s) %s-%s x%s", ph[3], ph[5], ph[7], ph[9])
	}
	if international.MatchString(in) {
		// remove spaces after +
		return regexp.MustCompile(`^+( +)`).ReplaceAllString(in, "")
	}
	return in
}

func validateAuthGrantsJS(v interface{}, _ cty.Path) diag.Diagnostics {
	js := []byte(v.(string))
	if len(js) == 0 {
		return nil
	}

	var authGrants []iam.AuthGrantRequest
	if err := json.Unmarshal(js, &authGrants); err != nil {
		return diag.Errorf("auth_grants_json is not valid: %s", err)
	}

	if len(authGrants) == 0 {
		return diag.Errorf("auth_grants_json must contain at least one entry")
	}

	return nil
}

// UnknownVariableValue is a sentinel value that is used
// to denote that the value of a variable is unknown at this time.
// RawConfig uses this information to build up data about
// unknown keys.
// https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.17.0/internal/configs/hcl2shim/values.go#L16
const UnknownVariableValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

func stateAuthGrantsJS(v interface{}) string {
	if v.(string) == UnknownVariableValue {
		return UnknownVariableValue
	}
	js := []byte(v.(string))
	if len(js) == 0 {
		return ""
	}

	var authGrants []iam.AuthGrantRequest
	if err := json.Unmarshal(js, &authGrants); err != nil {
		panic(fmt.Sprintf(`"auth_grants": %q is not valid: %s`, v.(string), err))
	}

	var authGrantsJSON []byte
	authGrantsJSON, err := json.Marshal(authGrants)
	if err != nil {
		panic(fmt.Sprintf(`"auth_grants": %q is not valid: %s`, v.(string), err))
	}

	return string(authGrantsJSON)
}

func suppressAuthGrantsJS(k, oldVal, newVal string, _ *schema.ResourceData) bool {
	if newVal == UnknownVariableValue {
		return false
	}

	var oldAuthGrants []iam.AuthGrantRequest
	if len(oldVal) > 0 {
		if err := json.Unmarshal([]byte(oldVal), &oldAuthGrants); err != nil {
			panic(fmt.Sprintf("previous value for %q: %q is not valid: %s", k, oldVal, err))
		}
	}

	var newAuthGrants []iam.AuthGrantRequest
	if len(newVal) > 0 {
		if err := json.Unmarshal([]byte(newVal), &newAuthGrants); err != nil {
			panic(fmt.Sprintf("new value for %q: %q is not valid: %s", k, newVal, err))
		}
	}

	// sort grants before comparing; swaping grants order should not cause update
	sortAuthGrants(oldAuthGrants)
	sortAuthGrants(newAuthGrants)

	return cmp.Equal(oldAuthGrants, newAuthGrants, cmpopts.EquateEmpty())
}

func sortAuthGrants(grants []iam.AuthGrantRequest) {
	sort.Slice(grants, func(i, j int) bool {
		return grants[i].GroupID < grants[j].GroupID
	})
	for _, g := range grants {
		sortAuthGrants(g.Subgroups)
	}
}

func statePhone(v interface{}) string {
	return canonicalPhone(v.(string))
}

func suppressPhone(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	oldVal = regexp.MustCompile(`\D+`).ReplaceAllLiteralString(oldVal, "")
	newVal = regexp.MustCompile(`\D+`).ReplaceAllLiteralString(newVal, "")
	return oldVal == newVal
}

func suppressEmail(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	return strings.EqualFold(oldVal, newVal)
}

func stateEmail(v interface{}) string {
	return strings.ToLower(v.(string))
}

func validatePhone(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be a string, but got %T", k, k))
		return warnings, errors
	}

	if v == "" {
		return warnings, errors
	}
	if !northAmerica.MatchString(v) && !international.MatchString(v) {
		errors = append(errors, fmt.Errorf("%q contains invalid phone number: %q", k, v))
	}

	return warnings, errors
}

func customizePasswordDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.HasChange("password") {
		oldVal, newVal := d.GetChange("password")
		oldPassword := oldVal.(string)
		newPassword := newVal.(string)
		if oldPassword != "" && newPassword == "" {
			return fmt.Errorf("deleting the password field or setting the password to an empty string is not allowed")
		}
	}
	return nil
}

func customizeNotificationDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	newValue := d.GetRawConfig()
	if newValue.GetAttr("user_notifications").LengthInt() == 0 {
		// Field is omitted, so apply the default configuration
		defaultConfig := []interface{}{
			map[string]interface{}{
				"enable_email_notifications":                true,
				"api_client_credential_expiry_notification": false,
				"new_user_notification":                     true,
				"password_expiry":                           true,
				"proactive":                                 []interface{}{},
				"upgrade":                                   []interface{}{},
			},
		}
		if err := d.SetNew("user_notifications", defaultConfig); err != nil {
			return fmt.Errorf("failed to set default notification configuration: %s", err)
		}
	}

	return nil
}
