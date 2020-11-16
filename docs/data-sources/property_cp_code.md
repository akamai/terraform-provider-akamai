---
layout: "akamai"
page_title: "Akamai: cp_code"
subcategory: "Provisioning"
description: |-
 CP Code
---

# akamai_cp_code


Use `akamai_cp_code` data source to retrieve a cpcode id.

## Example Usage

Basic usage:

```hcl
data "akamai_cp_code" "example" {
     name = "my cpcode name"
     group_id = "grp_123"
     contract_id = "ctr_1-AB123"
}
```

A more real world example using other datasources as dependencies:
```
locals {
    group_name = "example group name"
    cpcode_name = "My Cpcode Name"
}

data "akamai_group" "example" {
    name = local.group_name
    contract_id = data.akamai_contract.example.id
}

data "akamai_contract" "example" {
     group_name = local.group_name
}

data "akamai_cp_code" "example" {
     name = local.cpcode_name
     group_id = data.akamai_group.example.id
     contract_id = data.akamai_contract.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The CP code name.
* `group_id` — (Required) The group ID
* `contract_id` — (Required) The contract ID

### Deprecated Arguments
* `group` — (Deprecated) synonym of group_id for legacy purposes. Cannot be used with `group_id`
* `contract` — (Deprecated) synonym of contract_id for legacy purposes. Cannot be used with `contract_id`

## Attributes Reference

The following attributes are returned:

* `id` — The CP code ID.
* `product_ids` - An array of product ids associated with this cpcode
