---
layout: "akamai"
page_title: "Akamai: groups"
subcategory: "Common"
description: |-
 Property groups
---

# akamai_groups


Use the `akamai_property_groups` data source to list groups associated with the [EdgeGrid API client token](https://developer.akamai.com/getting-started/edgegrid) you're using.

## Basic Usage

Return groups associated with the selected EdgeGrid API client token:

datasource-example.tf
```hcl-terraform
datasource "akamai_groups" "my-example" {
}

output "property_match" {
  value = data.akamai_groups.my-example
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

This data source returns these attributes:

* `groups` - A list of supported groups, with the following properties:
  * `group_id` - A group's unique ID. If your ID doesn't include the `grp_` prefix, the Akamai Provider appends it to your entry for processing purposes.
  * `group_name` - The name of the group.
  * `parent_group_id` - The ID of the parent group, if applicable.
  * `contract_ids` - An array of strings listing the contract IDs for each group.