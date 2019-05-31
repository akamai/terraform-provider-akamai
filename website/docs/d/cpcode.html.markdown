---
layout: "akamai"
page_title: "Akamai: cpcode"
sidebar_current: "docs-akamai-data-cpcode"
description: |-
 CPCode
---

# akamai_cpcode


Use `akamai_cpcode` datasource to retrieve a group id.



## Example Usage

Basic usage:

```hcl
data "akamai_cpcode" "terraform-demo" {
     name = "cpcode_name"
}

resource "akamai_property" "terraform-demo-web" {
     contract = "${data.akamai_cpcode.terraform-demo.id}"
...
}




```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The cpcode.

## Attributes Reference

The following are the return attributes:

* `id` — The cpcode.
