---
layout: "akamai"
page_title: "Akamai: authorities_set"
subcategory: "DNS"
description: |-
 DNS Authorities Set
---

# akamai_authorities_set

Use `akamai_authorities_set` datasource to retrieve a contracts authorities set for use when creating new zones.

## Example Usage

Basic usage:

```hcl
data "akamai_authorities_set" "example" {
     contract = "ctr_1-AB123"
}
```

## Argument Reference

The following arguments are supported:

* `contract` — (Required) The contract ID.

## Attributes Reference

The following attributes are returned:

* `authorities` — A list of authorities
