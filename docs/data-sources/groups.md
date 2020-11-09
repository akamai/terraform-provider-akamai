---
layout: "akamai"
page_title: "Akamai: groups"
subcategory: "Common"
description: |-
 Group
---

# akamai_groups


Use `akamai_groups` data source to retrieve the list of supported groups.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
    edgerc = "~/.edgerc"
}

data "akamai_groups" "example" {
}

output "groups" {
    value = "${data.akamai_groups.example.groups}"
}
```

## Argument Reference

No arguments needed.

## Attributes Reference

Following are the return attributes:

* `groups` — list of supported groups, with the following properties:
  * `group_id` - the group ID (string)
  * `group_name` - the group name (string)
  * `parent_group_id` - the parent group ID (string)
  * `contract_ids` - the group contract IDs (array of string)