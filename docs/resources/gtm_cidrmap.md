---
layout: "akamai"
page_title: "Akamai: gtm cidrmap"
subcategory: "Global Traffic Management"  
description: |-
  GTM Cidr Map
---

# akamai_gtm_cidrmap

`akamai_gtm_cidrmap` provides the resource for creating, configuring and importing a GTM Classless Inter-Domain Routing (CIDR) Map. CIDR mapping uses the IP addresses of the requesting name server to provide IP-specific CNAME entries, which lets you direct internal users to a specific environment or direct to the origin. This provides different responses to an internal corporate DNS infrastructure such as internal test environments and another answer for all other (defaultDatacenter) name servers. CIDR maps split the Internet into multiple CIDR block zones. Properties that use a map can specify a handout CNAME for each zone on the property’s editing page. To configure a property for CIDR mapping, your domain needs at least one CIDR map defined.. Note: Import requires an ID of the format: `existing_domain_name`:`existing_map_name`.

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

### Required

* `domain` - GTM Domain name for the AS Map.
* `name` - A descriptive label for the CIDR map, up to 255 characters.
* `default_datacenter` - A placeholder for all other CIDR zones not found in these CIDR zones.
  * `datacenter_id` - For each property, an identifier for all other CIDR zones.
  * `nickname` - A descriptive label for the all other CIDR blocks.

### Optional
 
* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete.
* `assignment` - (multiple allowed) Contains information about the CIDR zone groupings of CIDR blocks.
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `nickname` - A descriptive label for the CIDR zone group, up to 256 characters.
  * `blocks` - (List) Specifies an array of CIDR blocks.

### Schema Reference

The GTM Cidr Map backing schema and more complete element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#cidrmap).

