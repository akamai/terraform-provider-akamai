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

This data source supports these arguments:

### Required Arguments
* Group qualifier in one of the three forms detailed below.  Used to keep group and contract selections in synch when using an API that requires both.
  * `group_name` - The group name within which the contract can be found. 
  * `group_id` - The group id within which the contract can be found. 
  * `group` - (Deprecated) Either a group id or a group name within which the contract can be found. Cannot be used with `group_id` and `group_name`.

## Attributes Reference

* `id` - The contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it to your entry for processing purposes.

