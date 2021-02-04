---
layout: "akamai"
page_title: "Akamai: gtm_default_datacenter"
subcategory: "Global Traffic Management"
description: |-
 Default data center
---

# akamai_gtm_default_datacenter

Use the `akamai_gtm_default_datacenter` data source to retrieve the default data center, ID, and nickname.

## Example usage

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

## Argument reference

This data source supports these arguments:

* `domain` - (Required)
* `datacenter` - (Optional) The default is `5400`.

## Attributes reference

This data source supports these attributes:

* `id` - The data resource ID. Enter in this format: `<domain>:default_datacenter:<datacenter_id>`.
* `datacenter_id` - The default data center ID.
* `nickname` - The default data center nickname.
