---
layout: "akamai"
page_title: "Akamai: gtm asmap"
subcategory: "Global Traffic Management"  
description: |-
  GTM AS Map
---

# akamai_gtm_asmap

`akamai_gtm_asmap` provides the resource for creating, configuring and importing a GTM Autonomous System Map. Autonomous System (AS) mapping lets you configure a GTM property that returns a CNAME based on the AS number of the IP address of the requester. You can reuse maps for multiple properties or create new ones. AS maps split the Internet into multiple AS block zones. Properties that use AS maps can specify handout integers for each zone. AS mapping lets you configure a property that directs users to a specific environment or to the origin. Note: Import requires an ID of the format: `existing_domain_name`:`existing_map_name`.

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

### Required

* `domain` - GTM Domain name for the AS Map.
* `name` - A descriptive label for the AS map. Properties set up for asmapping can use this as reference.
* `default_datacenter` - A placeholder for all other AS zones not found in these AS zones.
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `nickname` - A descriptive label for all other AS zones, up to 128 characters.

### Optional
 
* `wait_on_complete` - (Boolean, Default: `true`) Wait for transaction to complete.
* `assignment` - (multiple allowed) Contains information about the AS zone groupings of AS IDs.
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `nickname` - A descriptive label for the group.
  * `as_numbers` - (List) Specifies an array of AS numbers.

### Schema Reference

The GTM AS Map backing schema and more complete element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#asmap)

