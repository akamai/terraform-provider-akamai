---
layout: "akamai"
page_title: "Akamai: akamai_iam_user"
subcategory: "Identity and Access Management"
description: |-
  Create user resources.
---

# akamai_iam_user

The `akamai_iam_user` resource represents a user on the Akamai platform. 

## Basic usage

This example shows how to set up a user:

```hcl
resource "akamai_iam_user" "example_user" {
  first_name       = "Juan"
  last_name        = "Perez"
  phone            = "+1 206-555-0100"
  email            = "jperez@example.com"
  country          = "Grenada"
  enable_tfa       = false
  auth_grants_json = "[{\"groupId\":18451,\"roleId\":14},{\"groupId\":18453,\"roleId\":13}]"
}
```

## Argument reference

This resource supports these arguments:

* `first_name` - (Required) The user's first name.
* `last_name` - (Required) The user's last name.
* `email` - (Required) The user's email address.
* `country` - (Required) As part of the user's location, the value can be any that are available from the [view-supported-countries operation](../data-sources/iam_countries.md).
* `phone` - (Required) The user's main phone number.
* `enable_tfa` - (Required) Indicates whether two-factor authentication is allowed.
* `auth_grants_json` - (Required) A user's per-group role assignments, in JSON form.
* `contact_type` - (Optional) To help characterize the user, the value can be any that are available from the [view-contact-types operation](../data-sources/iam_contact_types.md).
* `job_title` - (Optional) The user's position at your company
* `time_zone` - (Optional) The user's time zone. The value can be any that are available from the [view-time-zones operation](../data-sources/iam_timezones.md)
* `secondary_email` - (Optional) The user's secondary email address.
* `mobile_phone` - (Optional) The user's mobile phone number.
* `address` - (Optional) The user's street address.
* `city` - (Optional) The user's city.
* `state` - (Optional) The user's state.
* `zip_code` - (Optional) The user's five-digit ZIP code.
* `preferred_language` - (Optional) The value can be any that are available from the [view-languages operation](../data-sources/iam_supported_langs.md)

### Deprecated arguments
* `is_locked` - (Deprecated) The user's lock status. 

## Attributes reference

This resource returns these attributes:

* `session_timeout` - The number of seconds it takes for the user's session to time out if there hasn't been any activity.
* `user_name` - A user's `loginId`. Typically, a user's email address.
* `is_locked` - The user's lock status.
* `last_login` - ISO 8601 timestamp indicating when the user last logged in.
* `password_expired_after` - The date a user's password expires.
* `tfa_configured` - Indicates whether two-factor authentication is configured.
* `email_update_pending` - Indicates whether email update is pending.
* `lock` - (Optional) Flag to block as user account.
