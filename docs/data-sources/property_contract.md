---
layout: "akamai"
page_title: "Akamai: contract"
subcategory: "Common"
description: |-
 Contract
---

# akamai_contract


Use the `akamai_contract` data source to find a contract ID.

## Example usage

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

## Argument reference

This data source requires one of these group arguments to return contract information: 
  * `group_name` - The name of the group containing the contract. 
  * `group_id` - The unique ID of the group containing the contract, including the  `grp_` prefix.

### Deprecated arguments

* `group` - (Deprecated) Either the group ID or the group name that includes the contract. You can't use this argument with `group_id` and `group_name`.

## Attributes reference

* `id` - The contract's unique ID, including the `ctr_` prefix.

