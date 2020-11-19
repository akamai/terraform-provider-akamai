---
layout: "akamai"
page_title: "Akamai: contract"
subcategory: "Common"
description: |-
 Contract
---

# akamai_contract


Use the `akamai_contract` data source to find a contract ID.

## Example Usage

Basic usage:

```hcl
data "akamai_contract" "example" {
     group_name = "example group name"
}

resource "akamai_property" "example" {
    contract_id = data.akamai_contract.example.id
    ...
}
```

## Argument Reference

This data source requires one of these group arguments to return contract information: 
  * `group_name` - The name of the group containing the contract. 
  * `group_id` - The unique ID of the group containing the contract. If your ID doesn't include the `grp_` prefix, the Akamai provider appends it to your entry for processing purposes.

### Deprecated Arguments

* `group` - (Deprecated) Either the group ID or the group name that includes the contract. You can't use this argument with `group_id` and `group_name`.

## Attributes Reference

* `id` - The contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it to your entry for processing purposes.

