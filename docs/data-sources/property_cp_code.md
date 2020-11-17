---
layout: "akamai"
page_title: "Akamai: cp_code"
subcategory: "Provisioning"
description: |-
 CP Code
---

# akamai_cp_code


Use the `akamai_cp_code` data source to retrieve the ID for a content provider (CP) code.

## Example Usage

Basic usage:

```hcl
data "akamai_cp_code" "example" {
     name = "my cpcode name"
     group_id = "grp_123"
     contract_id = "ctr_1-AB123"
}
```

Here's a more real-world example that includes other data sources as dependencies:
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

* `name` - (Required) The name of the CP code.
* `group_id` - The group's unique ID. If your ID doesn't include the `grp_` prefix, the Akamai Provider appends it to your entry for processing purposes.
* `contract_id` - (Required) A contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it to your entry for processing purposes. 

### Deprecated Arguments
* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `group` - (Deprecated) Replaced by `group_id`. Maintained for legacy purposes.

## Attributes Reference

This data source returns these attributes are returned:

* `id` - The CP code ID.
* `product_ids` - An array of product IDs associated with this CP code. Each ID will include the `prd_` prefix.
