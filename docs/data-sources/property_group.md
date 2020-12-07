---
layout: "akamai"
page_title: "Akamai: group"
subcategory: "Common"
description: |-
 Group
---

# akamai_group


Use the `akamai_group` data source to retrieve a group ID by name. 

Each account features a hierarchy of groups, which control access to your Akamai configurations and 
help consolidate reporting functions, typically mapping to an organizational hierarchy. Using either 
Control Center or the [Identity Management: User Administration API](https://developer.akamai.com/en-us/api/core_features/identity_management_user_admin/v2.html), 
account administrators can assign properties to specific groups, each with its own set of users and 
accompanying roles.

## Example Usage

Basic usage:

```hcl
data "akamai_group" "example" {
    name = "example group name"
    contract_id = data.akamai_contract.example.id
}

data "akamai_contract" "example" {
     group_name = "example group name"
}

resource "akamai_property" "example" {
    group_id    = data.akamai_group.example.id
    ...
}
```

## Argument Reference

This data source supports these arguments:

* `name` - (Required) The group name.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 

### Deprecated Arguments 
* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.

## Attributes Reference

This data source returns this attribute:

* `id` - The group's unique ID, including the `grp_` prefix.