---
layout: "akamai"
page_title: "Akamai: authorities_set"
sidebar_current: "docs-akamai-data-authorities-set"
description: |-
 DNS Authorities Set
---

# akamai_authorities_set

Use `akamai_authorities_set` datasource to retrieve a contracts authorities set for use when creating new zones.

## Example Usage

Basic usage:

```hcl
data "akamai_authorities_set" "example" {
     contract = "ctr_#####"
}
```

## Argument Reference

The following arguments are supported:

* `contract` — (Required) The contract ID.

## Attributes Reference

The following are the return attributes:

* `authorities` — A list of authorities
