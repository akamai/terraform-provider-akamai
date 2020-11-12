---
layout: "akamai"
page_title: "Akamai: group"
subcategory: "Common"
description: |-
 Group
---

# akamai_group


Use `akamai_group` data source to retrieve a group id by name.

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
    group_id    = ata.akamai_group.example.id
    ...
}
```

## Argument Reference

The following arguments are supported:

* `group_name` — (Required) The group name.
* `contract_id` — (Required) The contract ID. 

### Deprecated Arguments 
* `name` — (Deprecated) synonym for `group_name` for legacy purposes.
* `contract` — (Deprecated) synonym of contract_id for legacy purposes. 

## Attributes Reference

The following are the return attributes:

* `id` — The group ID.
