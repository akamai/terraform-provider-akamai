---
layout: "akamai"
page_title: "Akamai: akamai_iam_roles"
subcategory: "IAM"
description: |-
 IAM Roles
---

# akamai_iam_roles

Use `akamai_iam_roles` to list roles for the current account and contract type. The account and contract type are determined by the access tokens in your API client.

## Example Usage

Basic usage:

```hcl
data "akamai_iam_roles" "my-roles" {
    group_id = "1234567"
    get_actions = true
}

output "roles" {
  value = data.akamai_iam_roles.my-roles
}
```

## Argument Reference

The following arguments are supported:

* `group_id` — (optional, string) A unique identifier for a group.
* `get_actions` - (optional, bool) When enabled, the response includes information about actions such as "edit" or "delete"
* `get_users` - (optional, bool) When enabled, returns users assigned to the roles
* `ignore_context` - (optional, bool) When enabled, returns all roles for the current account without regard the contract type associated with your API client

## Attributes Reference

The following attributes are returned:

* `roles` — A list of roles

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getroles)