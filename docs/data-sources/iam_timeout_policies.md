---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_timeout_policies

Use `akamai_iam_timeout_policies` to list all the possible session timeout policies that Akamai supports. Use the values from this data source to set the session timeout for a user. The name for each timeout period is in minutes, and the time value is in seconds.

## Example usage

Basic usage:

```hcl
data "akamai_iam_timeout_policies" "timeout_policies" {
}

output "supported_timeout_policies" {
  value = data.akamai_iam_timeout_policies.timeout_policies
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `policies` — A map of session timeout policy names to their value in seconds.

[API Reference](https://techdocs.akamai.com/iam-api/reference/get-common-timeout-policies)
