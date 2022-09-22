---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_blocked_user_properties (Beta)

Use the `akamai_iam_blocked_user_properties` resource to remove or grant access to properties. Administrators can block a user's access to any property, overriding any available role already assigned to that user.

## Basic usage

This example returns the policy details based on the policy ID and optionally, a version:

```hcl
resource "akamai_iam_blocked_user_properties" "example" {
  identity_id        = "A-B-123456"
  group_id           = 12345
  blocked_properties = [1, 2, 3, 4, 5]
}
```

## Argument reference

This resource supports these arguments:

* `identity_id` - (Required) A unique identifier that corresponds to a user's actual profile or client ID. Each identifier must be a string.
* `group_id` - (Required) A unique identifier for a group. Each identifier must be an integer.
* `blocked_properties` - (Required) List of properties to block for a user. The property IDs must be an integer.


## Attributes reference

This resource doesn't return any attributes.
