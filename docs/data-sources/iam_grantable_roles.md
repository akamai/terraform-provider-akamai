---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_grantable_roles

List which grantable roles you can include in a new custom role or add to an existing custom role.

## Basic usage

This example returns the available roles to grant to users:

```hcl
data "akamai_iam_grantable_roles" "example" {}

output "aka_grantable_roles_count" {
  value = length(data.akamai_iam_grantable_roles.test.grantable_roles)
}

output "aka_grantable_roles" {
  value = data.akamai_iam_grantable_roles.test
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

This resource returns this attribute:

* `grantable_roles` - Lists which grantable roles you can include in a new custom role or add to an existing custom role.
  * `granted_role_id` - Granted role ID.
  * `name` - Granted role name.
  * `description` - Granted role description.