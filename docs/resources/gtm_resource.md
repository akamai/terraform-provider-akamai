---
layout: "akamai"
page_title: "Akamai: gtm resource"
subcategory: "Global Traffic Management"  
description: |-
  GTM Resource
---

# akamai_gtm_resource

`akamai_gtm_resource` provides the resource for creating, configuring and importing a GTM resource that represents a constraint on how much load a data center can absorb. Consider a Resource as something that can impose a capacity constraint on the load associated with one or more Properties in a Datacenter. Examples of Resources include: bandwidth, CPU load average, database queries per second, or disk operations per second. Note: Import requires an ID of the format: `existing_domain_name`:`existing_resource_name`.

## Example Usage

Basic usage:

```hcl
resource "akamai_gtm_resource" "demo_resource" {
    domain = "demo_domain.akadns.net"
    name = "demo_resource"
    aggregation_type = "latest"
    type = "XML load object via HTTP"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `domain` - DNS name for the GTM Domain set that includes this Property.
* `name` - A descriptive label for the Resource.
* `aggregation_type` - Specifies how GTM handles different load numbers when multiple load servers are used for a data center or property.
* `type` - Indicates the kind of loadObject format used to determine the load on the resource.

### Optional

* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete.
* `resource_instance`  - (multiple allowed) Contains information about the resources that constrain the properties within the data center.
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `load_object` - Identifies the load object file used to report real-time information about the current load, maximum allowable load, and target load on each resource.
  * `load_object_port` - Specifies the TCP port of the loadObject.
  * `load_servers` - (List) Specifies a list of servers from which to request the load object.
  * `use_default_load_object` - (Boolean) Indicates whether a default loadObject is used for the resources.
* `host_header` - Optionally specifies the host header used when fetching the load object.
* `least_squares_decay` - For internal use only. Unless Akamai indicates otherwise, omit the value or set it to null.
* `upper_bound` - An optional sanity check that specifies the maximum allowed value for any component of the load object.
* `description` - A descriptive note to help you track what the resource constrains.
* `leader_string` - Specifies the text that comes before the loadObject.
* `constrained_property` - Specifies the name of the property that this resource constrains, or ** to constrain all properties.
* `load_imbalance_percent` - Indicates the percent of load imbalance factor (LIF) for the property.
* `max_u_multiplicative_increment` - For internal use only. Unless Akamai indicates otherwise, omit the value or set it to null.
* `decay_rate` - For internal use only. Unless Akamai indicates otherwise, omit the value or set it to null.

### Schema Reference

The GTM Resource backing schema and more element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#resource)
