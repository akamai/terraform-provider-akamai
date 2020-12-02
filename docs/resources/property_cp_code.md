---
layout: "akamai"
page_title: "Akamai: CP Code"
subcategory: "Provisioning"
description: |-
  CP Code
---

# akamai_cp_code


The `akamai_cp_code` resource lets you create or reuse content provider (CP) codes.  CP codes track web traffic handled by Akamai servers. Akamai gives you a CP code when you purchase a product. You need this code when you activate associated properties. 

You can create additional CP codes to support more detailed billing and reporting functions.

By default, the Akamai Provider uses your existing CP code instead of creating a new one.

## Example Usage

Basic usage:

```hcl
resource "akamai_cp_code" "cp_code" {
  name = "My CP Code"
  contract_id = akamai_contract.contract.id
  group_id = akamai_group.group.id
  product_id = "prd_Object_Delivery"
}
```

Here's a real-life example that includes other data sources as dependencies:

```
locals {
    group_name = "example group name"
    cpcode_name = "My CP Code"
}

data "akamai_group" "example" {
    name = local.group_name
    contract_id = data.akamai_contract.example.id
}

data "akamai_contract" "example" {
    group_name = local.group_name
}

resource "akamai_cp_code" "example_cp" {
    name = local.cpcode_name
    group_id = data.akamai_group.example.id
    contract_id = data.akamai_contract.example.id
    product_id = "prd_Object_Delivery"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) A descriptive label for the CP code. If you're creating a new CP code, the name can’t include commas, underscores, quotes, or any of these special characters: ^ # %.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.
* `product_id` - (Required) A product's unique ID, including the `prd_` prefix.

### Deprecated Arguments

* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `group` - (Deprecated) Replaced by `group_id`. Maintained for legacy purposes.
* `product` - (Deprecated) Replaced by `product_id`. Maintained for legacy purposes.

## Attributes Reference

* `id` - The ID of the CP code.

## Import

Basic Usage:

```hcl
resource "akamai_cp_code" "example" {
    # (resource arguments)
  }
```

Akamai CP codes can be imported using a comma-delimited string of `cp_code_id,contract_id,group_id` in that order as ID, e.g.

```shell
$ terraform import akamai_cp_code.example cpc_123,ctr_1-AB123,grp_123
```