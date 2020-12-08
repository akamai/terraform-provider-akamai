---
layout: "akamai"
page_title: "Akamai: cp_code"
subcategory: "Provisioning"
description: |-
 CP Code
---

# akamai_cp_code


Use the `akamai_cp_code` data source to retrieve the ID for a content provider (CP) code.

## Example usage

Basic usage:

```hcl
data "akamai_cp_code" "example" {
     name = "my cpcode name"
     group_id = "grp_123"
     contract_id = "ctr_1-AB123"
}
```

Here's a real-world example that includes other data sources as dependencies:
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

## Argument reference

This data source supports these arguments:

* `name` - (Required) The name of the CP code.
* `group_id` - (Required) The group's unique ID, including the `grp_` prefix.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 

### Deprecated arguments
* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `group` - (Deprecated) Replaced by `group_id`. Maintained for legacy purposes.

## Attributes reference

This data source returns these attributes:

* `id` - The ID of the CP code, including the `cpc_` prefix.
* `product_ids` - An array of product IDs associated with this CP code. Each ID returned includes the `prd_` prefix.
