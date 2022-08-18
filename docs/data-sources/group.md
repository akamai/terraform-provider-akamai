---
layout: "akamai"
page_title: "Akamai: group"
subcategory: "Common"
description: |-
 Group
---

# akamai_group

Use the `akamai_group` data source to get a group by name.

Each account features a hierarchy of groups, which control access to your
Akamai configurations and help consolidate reporting functions, typically
mapping to an organizational hierarchy. Using either Control Center or the
[Identity Management: User Administration API](https://techdocs.akamai.com/iam-api/reference/api),
account administrators can assign properties to specific groups, each with
its own set of users and accompanying roles.

## Example usage

Basic usage:

```hcl
data "akamai_group" "example" {
    group_name = "example group name"
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

## Argument reference

This data source supports these arguments:

* `group_name` - (Required) The group name.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix.

### Deprecated arguments
* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `name` -  (Deprecated) Replaced by `group_name`. Maintained for legacy purposes.

## Attributes reference

This data source returns this attribute:

* `id` - The group's unique ID, including the `grp_` prefix.
