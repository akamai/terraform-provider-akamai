---
layout: "akamai"
page_title: "Akamai: authorities_set"
subcategory: "DNS"
description: |-
 DNS Authorities Set
---

# akamai_authorities_set

Use the `akamai_authorities_set` data source to retrieve a contract's authorities set. You use the authorities set when creating new zones.

## Example usage

Basic usage:

```
data "akamai_authorities_set" "example" {
     contract = "ctr_1-AB123"
}
```

## Argument reference

This data source supports this argument:

* `contract` - (Required) The contract ID.

## Attributes reference

This data source supports this attribute:

* `authorities` - A list of authorities.
