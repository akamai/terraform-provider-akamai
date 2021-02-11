---
layout: "akamai"
page_title: "Akamai: gtm geomap"
subcategory: "Global Traffic Management"  
description: |-
  GTM Geographic Map
---

# akamai_gtm_geomap

Use the `akamai_gtm_geomap` resource to create, configure, and import a GTM Geographic map. Geographic mapping lets you configure a property that returns a CNAME based on the geographic location of the request. 

You can reuse maps for multiple properties or create new ones. To configure a property for geographic mapping, you need to define at least one geographic map for your domain. Each map needs at least two definitions. For example, you can have one definition that maps a set of countries to a specific data center, and a second definition that routes all other traffic. 

~> **Note** Import requires an ID with this format: `existing_domain_name`:`existing_map_name`.

## Example usage

Basic usage:

```hcl
resource "akamai_gtm_geomap" "demo_geomap" {
    domain = "demo_domain.akadns.net"
    name = "demo_geo"
    default_datacenter {
        datacenter_id = 5400
        nickname = "All Others"
    }
}
```

## Argument reference

This resource supports these arguments:

* `domain` - (Required) GTM Domain name for the Geographic Map.
* `name` - (Required) A descriptive label for the Geographic map.
* `default_datacenter` - (Required) A placeholder for all other geographic zones. Requires these additional arguments:
  * `datacenter_id` - (Required) For each property, an identifier for all other geographic zones.
  * `nickname` - (Required) A descriptive label for all other geographic zones.
* `wait_on_complete` - (Optional) A boolean indicating whether to wait for transaction to complete. Set to `true` by default.
* `assignment` - (Optional) Contains information about the geographic zone groupings of countries. You can have multiple `assignment` arguments. If used, requires these additional arguments:
  * `datacenter_id` - (Optional) A unique identifier for an existing data center in the domain.
  * `nickname` - (Optional) A descriptive label for the group.
  * `countries` - (Optional) Specifies an array of two-letter ISO 3166 country codes, or for finer subdivisions, the two-letter country code and the two-letter stateOrProvince code separated by a forward slash.

## Schema reference

You can download the GTM Geographic Map backing schema from the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#geographicmap) page.
