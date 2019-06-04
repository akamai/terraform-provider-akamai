---
layout: "akamai"
page_title: "Akamai: contract"
sidebar_current: "docs-akamai-data-contract"
description: |-
 Contract
---

# akamai_contract


Use `akamai_contract` datasource to retrieve a group id.



## Example Usage

Basic usage:

```hcl
data "akamai_contract" "example" {
     name = "contract name"
}

resource "akamai_property" "example" {
    contract = "${data.akamai_contract.example.id}"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The contract name.

## Attributes Reference

The following are the return attributes:

* `id` — The contract ID.
