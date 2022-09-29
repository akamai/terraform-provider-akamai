---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_timezones

Use `akamai_iam_timezones` to list all time zones Akamai supports. Time zones are in ISO 8601 format. Use the values from this data source to set the time zone for a user. Administrators use this data source to set a user's time zone. The default time zone is GMT.

## Example usage

Basic usage:

```hcl
data "akamai_iam_timezones" "test" {
}

output "aka_timezone_count" {
  value = length(data.akamai_iam_timezones.test.timezones)
}

output "aka_timezones" {
  value = data.akamai_iam_timezones.test.timezones
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `timezones` — Supported timezones.
  * `timezone` - The time zone ID.
  * `description` - The description of a time zone, including the GMT +/-.
  * `offset` - The time zone offset from GMT.
  * `posix` - The time zone posix.

[API Reference](https://techdocs.akamai.com/iam-api/reference/get-common-timezones)
