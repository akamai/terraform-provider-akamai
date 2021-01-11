---
layout: "akamai"
page_title: "Akamai: akamai_iam_roles"
subcategory: "IAM"
description: |-
 IAM Roles
---

# akamai_iam_roles

Use `akamai_iam_roles` to list roles for the current account and contract type. The account and contract type are determined by the access tokens in your API client. Use the `roleId` from this data source to construct the `auth_grants_json` when creating or updating a user's auth grants.

## Example usage

Basic usage:

```hcl
data "akamai_iam_roles" "my-roles" {}

output "roles" {
  value = data.akamai_iam_roles.my-roles
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `roles` — A list of roles.

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getroles)
