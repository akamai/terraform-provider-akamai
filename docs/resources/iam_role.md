---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_role

Use the `akamai_iam_role` resource to list and create roles for users. Roles are lists of permissions that are explicitly tied to both a user and a group. Users need roles to act on objects in a group.

## Basic usage

This example returns information on available roles:

```hcl
resource "akamai_iam_role" "example" {
    name           = "View Only"
    description    = "This role will allow you to view"
    granted_roles  = 2051
    type           = "custom"
}
```

## Argument reference

This resource supports these arguments:

* `name` - (Required) The name you supply for a role.
* `description` - (Required) The description for a role.
* `granted_roles` - (Required) The list of existing unique identifiers for the granted roles. Each identifier must be a unique integer.

## Attributes reference

This resource returns this attribute:

* `type` - The type indicates whether the role is `standard`, provided by Akamai, or `custom` for the account.




