---
layout: "akamai"
page_title: "Akamai: gtm property"
sidebar_current: "docs-akamai-resource-gtm-property"
description: |-
  GTM Property
---

# akamai_gtm_property

`akamai_gtm_property` provides the resource for creating, configuring and importing a gtm property to integrate easily with your existing GTM infrastructure to provide a secure, high performance, highly available and scalable solution for Global Traffic Management. Note: Import requires an ID of the format: `existing_domain_name`:`existing_property_name`

## Example Usage

Basic usage:

```hcl
resource "akamai_gtm_property" "demo_property" {
    domain = "demo_domain.akadns.net"
    name = "demo_property"
    type =  "weighted-round-robin"
    score_aggregation_type = "median"
    handout_limit = 5
    handout_mode = "normal"
}
```

## Argument Reference

The following arguments are supported:

Required

* `domain` — Domain name 
* `name` — Property name  
* `type` — Property type  
* `score_aggregation_type`
* `handout_limit` 
* `handout_mode`  
* `traffic_targets` — (List)
  * `datacenter_id`
  * `enabled` — (Boolean)
  * `weight`
  * `servers` — (List)
  * `name` — Traffic target name
  * `handout_cname`
* `liveness_tests` — (List)
  * `name` — Liveness test name
  * `test_interval`
  * `test_object_protocol`
  * `test_timeout`
  * `answer_required` — (Boolean)
  * `disable_nonstandard_port_warning` — (Boolean)
  * `error_penalty`
  * `host_header`
  * `http_error3xx` — (Boolean)
  * `http_error4xx` — (Boolean)
  * `http_error5xx` — (Boolean)
  * `peer_certificate_verification` — (Boolean)
  * `recursion_requested` — (Boolean)
  * `request_string`
  * `resource_type`
  * `response_string`
  * `ssl_client_certificate`
  * `ssl_client_private_key`
  * `test_object`
  * `test_object_password`
  * `test_object_port`
  * `test_object_username`
  * `timeout_penalty`

Optional
 
* `wait_on_complete` — (Boolean, Default: true) Wait for transaction to complete
* `failover_delay`
* `failback_delay`
* `ipv6` — (Boolean)
* `stickiness_bonus_percent`
* `stickiness_bonus_constant`
* `health_threshold`
* `use_computed_targets` — (Boolean)
* `backup_ip`
* `balance_by_download_score` — (Boolean)
* `static_ttl`
* `unreachable_threshold`
* `health_multiplier`
* `dynamic_ttl`
* `max_unreachable_penalty`
* `map_name`
* `load_imbalance_percentage`
* `health_max`
* `cname`
* `comments`
* `ghost_demand_reporting`
* `mx_records` — (List)
  * `exchange`
  * `preference`

Computed

The following arguments will be found in terraform.tfstate and can be referenced throughout the configuration. The values can NOT be changed.

* `weighted_hash_bits_for_ipv4`
* `weighted_hash_bits_for_ipv6`

### Backing Schema Reference

The GTM Property backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#property)

