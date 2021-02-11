---
layout: "akamai"
page_title: "Akamai: gtm datacenter"
subcategory: "Global Traffic Management"
description: |-
  GTM Datacenter
---

# akamai_gtm_datacenter

Use the `akamai_gtm_datacenter` resource to create, configure, and import a GTM data center. A GTM data center represents a customer data center and is also known as a traffic target, a location containing many servers GTM can direct traffic to. 

GTM uses data centers to scale load balancing. For example, you might have data centers in both New York and Amsterdam and want to balance load between them. You can configure GTM to send US users to the New York data center and European users to the data center in Amsterdam. 

~> **Note** Import requires an ID with this format: `existing_domain_name`:`existing_datacenter_id`.

## Example usage

Basic usage:

```hcl
resource "akamai_gtm_datacenter" "demo_datacenter" {
    domain = "demo_domain.akadns.net"
    nickname = "demo_datacenter"
}
```

## Argument reference

This resource supports these arguments:

* `domain` - (Required) The GTM domain name for the data center.
* `wait_on_complete` - (Optional) A boolean, that if set to `true`, waits for transaction to complete.
* `nickname` - (Optional) A descriptive label for the data center.
* `default_load_object` - (Optional) Specifies the load reporting interface between you and the GTM system. If used, requires these additional arguments:
  * `load_object` - A load object is a file that provides real-time information about the current load, maximum allowable load, and target load on each resource.
  * `load_object_port` - Specifies the TCP port to connect to when requesting the load object.
  * `load_servers` - Specifies a list of servers to request the load object from.
* `city` - (Optional) The name of the city where the data center is located.
* `clone_of` - (Optional) Identifies the data center’s `datacenter_id` of which this data center is a clone.
* `cloud_server_targeting` - (Optional) A boolean indicating whether to balance load between two or more servers in a cloud environment.
* `cloud_server_host_header_override` - (Optional) A boolean that, if set to `true`, Akamai's liveness test agents use the Host header configured in the liveness test.
* `continent` - (Optional) A two-letter code that specifies the continent where the data center maps to.
* `country` - (Optional) A two-letter ISO 3166 country code that specifies the country where the data center maps to.
* `latitude` - (Optional) Specifies the geographical latitude of the data center’s position. See also longitude within this object.
* `longitude` - (Optional) Specifies the geographic longitude of the data center’s position. See also latitude within this object.
* `state_or_province` - (Optional) Specifies a two-letter ISO 3166 country code for the state or province where the data center is located.

### Computed arguments

This resource returns these computed arguments in the `terraform.tfstate` file:

* `datacenter_id` - A unique identifier for an existing data center in the domain.
* `ping_interval`
* `ping_packet_size`
* `score_penalty`
* `servermonitor_liveness_count`
* `servermonitor_load_count`
* `servermonitor_pool`
* `virtual` - A boolean indicating whether the data center is virtual or physical, the latter meaning the data center has an Akamai Network Agent installed, and its physical location (`latitude`, `longitude`) is fixed. Either `true` if virtual or `false` if physical.

## Schema reference

You can download the GTM Data Center backing schema from the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#datacenter) page.
