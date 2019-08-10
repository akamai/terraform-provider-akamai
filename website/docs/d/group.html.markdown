---
layout: "akamai"
page_title: "Akamai: group"
sidebar_current: "docs-akamai-data-group"
description: |-
 Group
---

# akamai_group


Use `akamai_group` data source to retrieve a group id.

## Example Usage

Basic usage:

```hcl
data "akamai_group" "example" {
    name = "group name"
}


resource "akamai_property" "example" {
    group    = "${data.akamai_group.example.id}"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The group name.
* `contract` — (Optional) The contract ID

## Attributes Reference

The following are the return attributes:

* `id` — The group ID.
