---
layout: "akamai"
page_title: "Akamai: gtm domain"
subcategory: "docs-akamai-resource-gtm-domain"
description: |-
  GTM Domain
---

# akamai_gtm_domain

`akamai_gtm_domain` provides the resource for creating, configuring and importing a gtm domain to integrate easily with your existing GTM infrastructure to provide a secure, high performance, highly available and scalable solution for Global Traffic Management. Note: Import requires an ID of the format: `existing_domain_name`

## Example Usage

Basic usage:

```hcl
resource "akamai_gtm_domain" "demodomain" {
    contract = "XXX"
    group = 100
    name = "demo.akadns.net"
    type =  "basic"
    comment =  "some comment"
}
```

## Argument Reference

The following arguments are supported:

Required

* `contract` — The contract ID (if creating domain) 
* `group` — The currently selected group ID (if creating domain)   
* `name` — Domain name  
* `type` — Domain type  

Optional 

* `wait_on_complete` — (Boolean, Default: true) Wait for transaction to complete
* `comment` — A descriptive comment
* `email_notification_list` — (List)
* `default_timeout_penalty` — (Default: 25)
* `load_imbalance_percentage`
* `default_ssl_client_private_key`
* `default_error_penalty` — (Default: 75)
* `cname_coalescing_enabled` — (Boolean)
* `load_feedback` — (Boolean)
* `default_ssl_client_certificate`
* `end_user_mapping_enabled` — (Boolean)

Computed

The following arguments will be found in terraform.tfstate and can be referenced throughout the configuration. The values can NOT be changed.

* `default_unreachable_threshold` 
* `min_pingable_region_fraction`
* `servermonitor_liveness_count`
* `round_robin_prefix`
* `servermonitor_load_count`
* `ping_interval`
* `max_ttl`
* `default_health_max`
* `map_update_interval`
* `max_properties`
* `max_resources`
* `default_error_penalty`
* `max_test_timeout`
* `default_health_multiplier`
* `servermonitor_pool`
* `min_ttl`
* `default_max_unreachable_penalty`
* `default_health_threshold`
* `min_test_interval`
* `ping_packet_size`

### Backing Schema Reference

The GTM Domain backing schema and element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#domain)

