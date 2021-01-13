---
layout: "akamai"
page_title: "Akamai: groups"
subcategory: "Common"
description: |-
 Property groups
---

# akamai_groups


Use the `akamai_property_groups` data source to list groups associated with the [EdgeGrid API client token](https://developer.akamai.com/getting-started/edgegrid) you're using.

## Basic usage

Return groups associated with the EdgeGrid API client token you're using:

```hcl
data "akamai_groups" "my-example" {
}

output "property_match" {
  value = data.akamai_groups.my-example
}
```

## Argument reference

There are no arguments available for this data source.

## Attributes reference

This data source returns these attributes:

* `groups` - A list of supported groups, with the following attributes:
  * `group_id` - A group's unique ID, including the `grp_` prefix.
  * `group_name` - The name of the group.
  * `parent_group_id` - The ID of the parent group, if applicable.
  * `contract_ids` - An array of strings listing the contract IDs for each group.