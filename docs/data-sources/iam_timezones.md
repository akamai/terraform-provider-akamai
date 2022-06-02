---
layout: "akamai"
page_title: "Akamai: akamai_iam_timezones"
subcategory: "Identity and Access Management"
description: |-
 IAM Timeout Policies
---

# akamai_iam_timezones

Use `akamai_iam_timezones` to list all time zones Akamai supports. Time zones are in ISO 8601 format. Use the values from this data source to set the time zone for a user. Administrators use this data source to set a user's time zone. The default time zone is GMT.

## Example usage

Basic usage:

```hcl
data "akamai_iam_timezones" "timezones" {
}

output "supported_timezones" {
  value = data.akamai_iam_timezones.timezones
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `timezone` — The time zone ID.
* `description` - The description of a time zone, including the GMT +/-.
* `offset` - The time zone offset from GMT.
* `posix` - The time zone posix.

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getadmintimezones)
