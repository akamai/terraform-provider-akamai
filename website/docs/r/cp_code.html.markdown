---
layout: "akamai"
page_title: "Akamai: CP Code"
sidebar_current: "docs-akamai-resource-cp-code"
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
  contract = "${akamai_contract.contract.id}"
  group = "${akamai_group.group.id}"
  product = "prd_xxx"
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The CP Code name
* `contract` — (Required) The Contract ID
* `group` — (Required) The Group ID
* `product` — (Required) The Product ID

## Import Command

Import into Terraform state using CP Code ID, Contract ID and Group ID:

```
terraform import akamai_cp_code.foo cpc_123456:ctr_123456:grp_123456
```
All three IDs are required and must be in order delimeted with colons.