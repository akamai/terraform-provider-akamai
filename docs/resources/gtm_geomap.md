---
layout: "akamai"
page_title: "Akamai: gtm geomap"
subcategory: "Global Traffic Management"  
description: |-
  GTM Geographic Map
---

# akamai_gtm_geomap

`akamai_gtm_geomap` provides the resource for creating, configuring and importing a GTM Geographic map. Geographic mapping lets you configure a property that returns a CNAME based on the geographic location of the request. You can reuse maps for multiple properties or create new ones. To configure a property for geographic mapping, your domain needs at least one geographic map defined. Each map needs at least two definitions. This ensures that at least one definition maps one or more countries to a data center, and the second definition routes all other traffic. Note: Import requires an ID of the format: `existing_domain_name`:`existing_map_name`.

## Example Usage

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

## Argument Reference

The following arguments are supported:

### Required

* `domain` - GTM Domain name for the Geographic Map.
* `name` - A descriptive label for the Geographic map.
* `default_datacenter` - A placeholder for all other geographic zones.
  * `datacenter_id` - For each property, an identifier for all other geographic zones.
  * `nickname` - A descriptive label for all other geographic zones.

### Optional
 
* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete.
* `assignment` - (multiple allowed) Contains information about the geographic zone groupings of countries.
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `nickname` - A descriptive label for the group.
  * `countries` - (List) Specifies an array of two-letter ISO 3166 country codes, or for finer subdivisions, the two-letter country code and the two-letter stateOrProvince code separated by a forward slash.

### Schema Reference

The GTM Geographic Map backing schema and more complete element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#geographicmap).
