---
layout: "akamai"
page_title: "Akamai: gtm_default_datacenter"
sidebar_current: "docs-akamai-data-gtm_default_datacenter"
description: |-
 CP Code
---

# akamai_gtm_default_datacenter

Use `akamai_gtm_default_datacenter` data source to retrieve default datacenter id and nickname.

## Example Usage

Basic usage:

```hcl
data "akamai_gtm_default_datacenter" "example_ddc" {
     name = "example_domain.akadns.net"
     datacenter = 5400
}

resource "akamai_gtm_cidrmap" "example_cidrmap" {
    domain = "example_domain.akadns.net"
    default_datacenter {
        datacenter_id = data.akamai_gtm_default_datacenter.example.datacenter_id
        nickname = data.akamai_gtm_default_datacenter.example.nickname
    ...
}
```

## Argument Reference

The following arguments are supported:

* `domain` — (Required)
* `domain` — (Optional. Default 5400)
* `datacenter_id` — (Computed - do not configure) default datacenter Id
* `nickname` — (Computed - do not configure) default datacenter nickname

## Attributes Reference

The following are the return attributes:

* `id` — The data resource id. Format: <domain>:default_datacenter:<datacenter_id>
