---
layout: "akamai"
page_title: "Akamai: gtm asmap"
subcategory: "Global Traffic Management"  
description: |-
  GTM AS Map
---

# akamai_gtm_asmap

Use the `akamai_gtm_asmap` resource to create, configure, and import a GTM Autonomous System (AS) map. AS mapping lets you configure a GTM property that returns a CNAME based on the AS number associated with the requester's IP address. 

You can reuse maps for multiple properties or create new ones. AS maps split the Internet into multiple AS block zones. Properties that use AS maps can specify handout integers for each zone. AS mapping lets you configure a property that directs users to a specific environment or to the origin. 

~> **Note** Import requires an ID with this format: `existing_domain_name`:`existing_map_name`.

## Example usage

Basic usage:

```
resource "akamai_gtm_asmap" "demo_asmap" {
    domain = "demo_domain.akadns.net"
    name = "demo_as"
    default_datacenter { 
        datacenter_id = 5400
        nickname = "All Other AS numbers"
    }
}

```

## Argument reference

This resource supports these arguments:

* `domain` - (Required) The GTM Domain name for the AS map.
* `name` - (Required) A descriptive label for the AS map. Properties set up for  AS mapping can use this as reference.
* `default_datacenter` - (Required) A placeholder for all other AS zones not found in these AS zones. Requires these additional arguments:
  * `datacenter_id` - (Required) A unique identifier for an existing data center in the domain.
  * `nickname` - (Required) A descriptive label for all other AS zones, up to 128 characters.
* `wait_on_complete` - (Optional) A boolean that, if `true`, waits for transaction to complete.
* `assignment` - (Optional) Contains information about the AS zone groupings of AS IDs. You can have multiple entries with this argument. If used, requires these arguments:
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `nickname` - A descriptive label for the group.
  * `as_numbers` - Specifies an array of AS numbers.

## Schema reference

You can download the GTM AS Map backing schema from the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#asmap) page.

