---
layout: "akamai"
page_title: "Akamai: gtm_default_datacenter"
subcategory: "Global Traffic Management"
description: |-
 Default Datacenter
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

* `domain` - (Required)
* `datacenter` - (Optional. Default 5400)

## Attributes Reference

The following attributes are returned:

* `id` - The data resource ID. Format: `<domain>:default_datacenter:<datacenter_id>`
* `datacenter_id` - The default datacenter ID
* `nickname` - The default datacenter nickname
