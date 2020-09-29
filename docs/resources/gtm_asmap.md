---
layout: "akamai"
page_title: "Akamai: gtm asmap"
subcategory: "docs-akamai-resource-gtm-asmap"  
description: |-
  GTM AS Map
---

# akamai_gtm_asmap

`akamai_gtm_asmap` provides the resource for creating, configuring and importing a gtm AS Map to integrate easily with your existing GTM infrastructure to provide a secure, high performance, highly available and scalable solution for Global Traffic Management. Note: Import requires an ID of the format: `existing_domain_name`:`existing_map_name`

## Example Usage

Basic usage:

```hcl
resource "akamai_gtm_asmap" "demo_asmap" {
    domain = "demo_domain.akadns.net"
    name = "demo_as"
    default_datacenter { 
        datacenter_id = 5400
        nickname = "All Other AS numbers"
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
  * `as_numbers` — (List)

### Backing Schema Reference

The GTM AS Map backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#asmap)

