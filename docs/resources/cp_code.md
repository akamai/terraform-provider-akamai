---
layout: "akamai"
page_title: "Akamai: CP Code"
subcategory: "Provisioning"
description: |-
  CP Code
---

# akamai_cp_code


The `akamai_cp_code` resource allows you to create or re-use CP Codes.

If the CP Code already exists it will be used instead of creating a new one.

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

A more real example using other datasources as dependencies:
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

* `name` — (Required) The CP Code name
* `contract_id` — (Required) The Contract ID
* `group_id` — (Required) The Group ID
* `product_id` — (Required) The Product ID

### Deprecated
* `contract` — (Deprecated) synonym of contract_id for legacy purposes
* `group` — (Deprecated) synonym of group_id for legacy purposes
* `product` — (Deprecated) synonym of product_id for legacy purposes

## Attributes Reference

* `id` — The CP code ID.
