---
layout: "akamai"
page_title: "Akamai: groups"
subcategory: "Common"
description: |-
 Property groups
---

# akamai_groups


Use `akamai_groups` data source to list groups associated with an EdgeGrid API client token. 

## Basic Usage

Return groups associated with the EdgeGrid API client token:

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

The following attributes are returned:

* `groups` — list of supported groups, with the following properties:
  * `group_id` - the group ID (string)
  * `group_name` - the group name (string)
  * `parent_group_id` - the parent group ID (string)
  * `contract_ids` - the group contract IDs (array of string)