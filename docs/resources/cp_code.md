---
layout: "akamai"
page_title: "Akamai: CP Code"
subcategory: "Property Provisioning"
description: |-
  CP Code
---

# akamai_cp_code

~> **Note** Version 1.0.0 of the Akamai Terraform Provider is now available for the Property Provisioning module. To upgrade to the new version, you have to update this resource. See [Upgrade to Version 1.0.0](../guides/1.0_migration.md) for details.

The `akamai_cp_code` resource lets you create or reuse content provider (CP) codes.  CP codes track web traffic handled by Akamai servers. Akamai gives you a CP code when you purchase a product. You need this code when you activate associated properties.

You can create additional CP codes to support more detailed billing and reporting functions.

By default, the Akamai Provider uses your existing CP code instead of creating a new one.

## Example usage

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
## Argument reference

The following arguments are supported:

* `name` - (Required) A descriptive label for the CP code. If you're creating a new CP code, the name can't include commas, underscores, quotes, or any of these special characters: ^ # %.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.
* `product_id` - (Required) A product's unique ID, including the `prd_` prefix. See [Common Product IDs](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/appendix#common-product-ids) for more information.

### Deprecated arguments

* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `group` - (Deprecated) Replaced by `group_id`. Maintained for legacy purposes.
* `product` - (Deprecated) Replaced by `product_id`. Maintained for legacy purposes.

## Attributes reference

* `id` - The ID of the CP code.

## Import

Basic Usage:

```hcl
resource "akamai_cp_code" "example" {
    # (resource arguments)
  }
```

You can import your Akamai CP codes using a comma-delimited string of the CP code,
contract, and group IDs. You have to enter the IDs in this order:

`cpcode_id,contract_id,group_id`

For example:

```shell
$ terraform import akamai_cp_code.example cpc_123,ctr_1-AB123,grp_123
```
