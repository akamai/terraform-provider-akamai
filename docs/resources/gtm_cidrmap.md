---
layout: akamai
subcategory: Global Traffic Management  
---

# akamai_gtm_cidrmap

Use the `akamai_gtm_cidrmap` resource to create, configure, and import a GTM Classless Inter-Domain Routing (CIDR) map. CIDR mapping uses the IP addresses of the requesting name server to provide IP-specific CNAME entries. CNAMEs let you direct internal users to a specific environment or direct them to the origin. This lets you provide different responses to an internal corporate DNS infrastructure, such as internal test environments and another answer for all other name servers (`default_datacenter`).

 CIDR maps split the Internet into multiple CIDR block zones. Properties that use a map can specify a handout CNAME for each zone on the property's editing page. To configure a property for CIDR mapping, your domain needs at least one CIDR map defined.

~> **Note** Import requires an ID with this format: `existing_domain_name`:`existing_map_name`.

## Example usage

Basic usage:

```
resource "akamai_gtm_cidrmap" "demo_cidrmap" {
    domain = "demo_domain.akadns.net"
    name = "demo_cidr"
    default_datacenter {
        datacenter_id = 5400
        nickname = "All Other CIDR Blocks"
    }
}
```

## Argument reference

This resource supports these arguments:

* `domain` - (Required) GTM Domain name for the CIDR Map.
* `name` - (Required) A descriptive label for the CIDR map, up to 255 characters.
* `default_datacenter` - (Required) A placeholder for all other CIDR zones not found in these CIDR zones. Requires these additional arguments:
  * `datacenter_id` - (Required) For each property, an identifier for all other CIDR zones.
  * `nickname` - (Required) A descriptive label for the all other CIDR blocks.
* `wait_on_complete` - (Optional) A boolean that, if set to `true`, waits for transaction to complete.
* `assignment` - (Optional) Contains information about the CIDR zone groupings of CIDR blocks. You can have multiple entries with this argument. If used, requires these additional arguments:
  * `datacenter_id` - (Optional) A unique identifier for an existing data center in the domain.
  * `nickname` - (Optional) A descriptive label for the CIDR zone group, up to 256 characters.
  * `blocks` - (Optional, list) Specifies an array of CIDR blocks.
