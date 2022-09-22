---
layout: akamai
subcategory: Global Traffic Management  
---

# akamai_gtm_resource

The `akamai_gtm_resource` lets you create, configure, and import a GTM resource. In GTM, a resource is anything you can measure whose scarcity affects load balancing. Examples of resources include bandwidth, CPU load average, database queries per second, or disk operations per second.

~> **Note** Import requires an ID with this format: `existing_domain_name`:
`existing_resource_name`.

## Example usage

Basic usage:

```
resource "akamai_gtm_resource" "demo_resource" {
    domain = "demo_domain.akadns.net"
    name = "demo_resource"
    aggregation_type = "latest"
    type = "XML load object via HTTP"
}
```

## Argument reference

This resource supports these arguments:

* `domain` - (Required) DNS name for the GTM Domain set that includes this property.
* `name` - (Required) A descriptive label for the GTM resource.
* `aggregation_type` - (Required) Specifies how GTM handles different load numbers when multiple load servers are used for a data center or property.
* `type` - (Required) Indicates the kind of `load_object` format used to determine the load on the resource.
* `wait_on_complete` - (Optional) A boolean indicating whether to wait for transaction to complete. Set to `true` by default.
* `resource_instance`  - (Optional) (multiple allowed) Contains information about the resources that constrain the properties within the data center. You can have multiple `resource_instance` entries. Requires these arguments:
  * `datacenter_id` - (Optional) A unique identifier for an existing data center in the domain.
  * `load_object` - (Optional) Identifies the load object file used to report real-time information about the current load, maximum allowable load, and target load on each resource.
  * `load_object_port` - (Optional) Specifies the TCP port of the `load_object`.
  * `load_servers` - (Optional) (List) Specifies a list of servers from which to request the load object.
  * `use_default_load_object` - (Optional) A boolean that indicates whether a default `load_object` is used for the resources.
* `host_header` - (Optional) Optionally specifies the host header used when fetching the load object.
* `least_squares_decay` - (Optional) For internal use only. Unless Akamai indicates otherwise, omit the value or set it to null.
* `upper_bound` - (Optional) An optional sanity check that specifies the maximum allowed value for any component of the load object.
* `description` - (Optional) A descriptive note to help you track what the resource constrains.
* `leader_string` - (Optional) Specifies the text that comes before the `load_object`.
* `constrained_property` - (Optional) Specifies the name of the property that this resource constrains, enter `**` to constrain all properties.
* `load_imbalance_percent` - (Optional) Indicates the percent of load imbalance factor (LIF) for the property.
* `max_u_multiplicative_increment` - (Optional) For Akamai internal use only. You can omit the value or set it to `null`.
* `decay_rate` - (Optional) For Akamai internal use only. You can omit the value or set it to `null`.
