---
layout: "akamai"
page_title: "Akamai: gtm resource"
subcategory: "docs-akamai-resource-gtm-resource"  
description: |-
  GTM Resource
---

# akamai_gtm_resource

`akamai_gtm_resource` provides the resource for creating, configuring and importing a gtm resource to integrate easily with your existing GTM infrastructure to provide a secure, high performance, highly available and scalable solution for Global Traffic Management. Note: Import requires an ID of the format: `existing_domain_name`:`existing_resource_name`

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

Required

* `domain` — Domain name 
* `name` — Resource name
* `aggregation_type`
* `type` — Resource type

Optional
 
* `wait_on_complete` — (Boolean, Default: true) Wait for transaction to complete
* `resource_instance`  — (multiple allowed) 
  * `datacenter_id`
  * `load_object`
  * `load_object_port`
  * `load_servers` — (List)
  * `use_default_load_object` — (Boolean)
* `host_header`
* `least_squares_decay`
* `upper_bound`
* `description`
* `leader_string`
* `constrained_property`
* `load_imbalance_percent`
* `max_u_multiplicative_increment`
* `decay_rate`

### Backing Schema Reference

The GTM Resource backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#resource)

