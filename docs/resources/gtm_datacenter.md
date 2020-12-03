---
layout: "akamai"
page_title: "Akamai: gtm datacenter"
subcategory: "Global Traffic Management"
description: |-
  GTM Datacenter
---

# akamai_gtm_datacenter

`akamai_gtm_datacenter` provides the resource for creating, configuring and importing a gtm datacenter to integrate easily with your existing GTM infrastructure to provide a secure, high performance, highly available and scalable solution for Global Traffic Management. Note: Import requires an ID of the format: `existing_domain_name`:`existing_datacenter_id`

## Example Usage

Basic usage:

```hcl
resource "akamai_gtm_datacenter" "demo_datacenter" {
    domain = "demo_domain.akadns.net"
    nickname = "demo_datacenter"
}
```

## Argument Reference

The following arguments are supported:

Required

* `domain` - Domain name 

Optional
 
* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete
* `nickname` - datacenter nickname
* `default_load_object`
  * `load_object`
  * `load_object_port`
  * `load_servers` - (List)
* `city`
* `clone_of`
* `cloud_server_targeting` - (Boolean)
* `cloud_server_host_header_override` - (Boolean)
* `continent`
* `country`
* `latitude`
* `longitude`
* `state_or_province`

Computed

The following arguments will be found in terraform.tfstate and can be referenced throughout the configuration. The values can NOT be changed.

* `datacenter_id`
* `ping_interval`
* `ping_packet_size`
* `score_penalty`
* `servermonitor_liveness_count`
* `servermonitor_load_count`
* `servermonitor_pool`
* `virtual` - (Boolean)

### Backing Schema Reference

The GTM Datacenter backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#datacenter)

