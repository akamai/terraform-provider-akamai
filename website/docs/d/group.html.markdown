---
layout: "akamai"
page_title: "Akamai: group"
sidebar_current: "docs-akamai-data-group"
description: |-
 Group
---

# akamai_group


Use `akamai_group` datasource to retrieve a group id.



## Example Usage

Basic usage:

```hcl
data "akamai_group" "terraform-demo" {
    name = "group name"
}


resource "akamai_property" "terraform-demo-web" {
    group    = "${data.akamai_group.terraform-demo.id}"
....
}



```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The group name.

## Attributes Reference

The following are the return attributes:

* `id` — The group ID.
