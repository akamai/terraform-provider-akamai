---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_group

Use the `akamai_iam_group` resource to list details about groups. Groups are organizational containers for the objects you use.  Groups can contain other groups, primary objects like properties, and secondary objects like [edge hostnames](../resources/edge_hostname.md) or [content provider (CP) codes](../resources/cp_code.md).

## Basic usage

This example returns the policy details based on the policy ID and optionally, a version:

```hcl
resource "akamai_iam_group" "example" {
  parent_group_id = 12345
  name            = "MyParentGroup"
}
```

## Argument reference

This resource supports these arguments:

* `parent_group_id` - (Required) A unique identifier for the parent group. Each identifier must be an integer.
* `name` - (Required) Human readable name for a group.


## Attributes reference

This resource returns this attribute:

* `sub_groups` - Sub-groups that are related to this group. Each identifier must be an integer.
