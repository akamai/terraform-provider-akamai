---
layout: "akamai"
page_title: "Akamai: gtm datacenter"
subcategory: "Global Traffic Management"
description: |-
  GTM Datacenter
---

# akamai_gtm_datacenter

`akamai_gtm_datacenter` provides the resource for creating, configuring and importing a GTM datacenter that represents a customer data center, also known as a traffic target, which is a location that contains many servers to which GTM directs traffic. More generally, a data center is a name for a set of possible answers that GTM can return for a query and the unit GTM uses to scale load balancing. For example, you might have data centers in New York and in Amsterdam, and want to balance load between them, but prefer that U.S. users are sent to the New York data center and European users are sent to the Amsterdam data center. Note: Import requires an ID of the format: `existing_domain_name`:`existing_datacenter_id`.

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

* `domain` - GTM Domain name for the Datacenter.

Optional
 
* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete.
* `nickname` - A descriptive label for the data center.
* `default_load_object` - Specifies the load reporting interface between you and the GTM system.
  * `load_object` - A load object is a file that provides real-time information about the current load, maximum allowable load, and target load on each resource.
  * `load_object_port` - Specifies the TCP port to connect to when requesting the load object.
  * `load_servers` - (List) Specifies a list of servers to request the load object from.
* `city` - The name of the city where the data center is located.
* `clone_of` - Identifies the data center’s datacenterId of which this data center is a clone.
* `cloud_server_targeting` - (Boolean) Balances load between two or more servers in a cloud environment.
* `cloud_server_host_header_override` - (Boolean) Akamai's liveness test agents will populate the Host header with the host header value configured in the liveness test.
* `continent` - A two-letter code that specifies the continent where the data center maps to.
* `country` - A two-letter ISO 3166 country code that specifies the country where the data center maps to.
* `latitude` - Specifies the geographical latitude of the data center’s position. See also longitude within this object.
* `longitude` - Specifies the geographic longitude of the data center’s position. See also latitude within this object.
* `state_or_province` - Specifies a two-letter ISO 3166 country code for the state or province where the data center is located.

Computed

The following arguments will be found in `terraform.tfstate` and can be referenced throughout the configuration. The values cannot be changed.

* `datacenter_id`.
* `ping_interval`.
* `ping_packet_size`.
* `score_penalty`.
* `servermonitor_liveness_count`.
* `servermonitor_load_count`.
* `servermonitor_pool`.
* `virtual` - (Boolean).

### Backing Schema Reference

The GTM Datacenter backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#datacenter).
