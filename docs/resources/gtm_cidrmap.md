---
layout: "akamai"
page_title: "Akamai: gtm cidrmap"
subcategory: "docs-akamai-resource-gtm-cidrmap"  
description: |-
  GTM Cidr Map
---

# akamai_gtm_cidrmap

`akamai_gtm_cidrmap` provides the resource for creating, configuring and importing a gtm Cidr Map to integrate easily with your existing GTM infrastructure to provide a secure, high performance, highly available and scalable solution for Global Traffic Management. Note: Import requires an ID of the format: `existing_domain_name`:`existing_map_name`

## Example Usage

Basic usage:

```hcl
resource "akamai_gtm_cidrmap" "demo_cidrmap" {
    domain = "demo_domain.akadns.net"
    name = "demo_cidr"
    default_datacenter {
        datacenter_id = 5400
        nickname = "All Other CIDR Blocks"
    }
}
```

## Argument Reference

The following arguments are supported:

Required

* `domain` — Domain name 
* `name` — Resource name
* `default_datacenter`
  * `datacenter_id`
  * `nickname`

Optional
 
* `wait_on_complete` — (Boolean, Default: true) Wait for transaction to complete
* `assignment` — (multiple allowed)
  * `datacenter_id`
  * `nickname`
  * `blocks` — (List)

### Backing Schema Reference

The GTM Cidr Map backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#cidrmap)

