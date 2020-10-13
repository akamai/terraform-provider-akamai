---
layout: "akamai"
page_title: "Akamai: contract"
subcategory: "Common"
description: |-
 Contract
---

# akamai_contract


Use `akamai_contract` data source to retrieve a contract id.

## Example Usage

Basic usage:

```hcl
data "akamai_contract" "example" {
     group = "group name"
}

resource "akamai_property" "example" {
    contract = "${data.akamai_contract.example.id}"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `group` — (Optional) The group within which the contract can be found. Used to keep group and contract selections in synch when using an API that requires both.

## Attributes Reference

The following are the return attributes:

* `id` — The contract ID.
