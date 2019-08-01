---
layout: "akamai"
page_title: "Akamai: cp_code"
sidebar_current: "docs-akamai-data-cpcode"
description: |-
 CP Code
---

# akamai_cp_code

Use `akamai_cp_code` datasource to retrieve a group id.

## Example Usage

Basic usage:

```hcl
data "akamai_cp_code" "example" {
     name = "cpcode name"
     group = "grp_#####"
     contract = "ctr_#####"
}

resource "akamai_property" "example" {
    contract = "${data.akamai_cpcode.example.id}"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The CP code name.
* `group` — (Required) The group ID
* `contract` — (Required) The contract ID

## Attributes Reference

The following are the return attributes:

* `id` — The CP code ID.
