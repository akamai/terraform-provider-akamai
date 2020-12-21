---
layout: "akamai"
page_title: "Akamai: user"
subcategory: "IAM"
description: |-
  Create user resources.
---

# akamai_iam_user

## Argument reference

This resource supports these arguments:

* `first_name` - (Required) The user's first name.
* `last_name` - (Required) The user's last name.
* `email` - (Required) The user's email address.
* `country` - (Required) As part of the user's location, the value can be any that are available from the [view-supported-countries operation](../data-sources/iam_countries.md).
* `phone` - (Required) The user's main phone number.
* `enable_tfa` - (Required) Indicates whether two-factor authentication is allowed.
* `send_otp_email` - (Required) Whether to send a one-time password to the newly-created user by email.
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
* `session_timeout` - (Optional) The number of seconds it takes for the user's Control Center session to time out if there hasn't been any activity.
* `enable_notifications` - (Optional) Whether to allow email notifications (notifications emails suspended unless `true`).
* `subscribe_new_users` - (Optional) Whether to send emails to group administrators when new users are created.
* `subscribe_password_expiration` - (Optional) Whether to send emails regarding password expiration.
* `subscribe_product_issues` - (Optional) Products for which the user receives notification emails about service issues.
* `subscribe_product_upgrades` - (Optional) Products for which the user receives notification emails about upgrades.
