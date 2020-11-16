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
    group_id    = data.akamai_group.example.id
    ...
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The group name.
* `contract_id` — (Required) The contract ID. 

### Deprecated Arguments 
* `contract` — (Deprecated) synonym of contract_id for legacy purposes. Cannot be used with `contract_id`.

## Attributes Reference

The following attributes are returned:

* `id` — The group ID.
