---
layout: "akamai"
page_title: "Akamai: gtm domain"
subcategory: "Global Traffic Management"
description: |-
  GTM Domain
---

# akamai_gtm_domain

`akamai_gtm_domain` provides the resource for creating, configuring and importing a GTM Domain, a basic building block of a traffic management configuration. Note: Import requires an ID of the format: `existing_domain_name`.

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

### Required

* `contract` - The contract ID (if creating domain).
* `group` - The currently selected group ID (if creating domain).
* `name` - DNS name for a collection of GTM Properties.
* `type` - GTM Domain type of `failover-only`, `static`, `weighted`, `basic` or `full`. 

### Optional 

* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete.
* `comment` - A descriptive note about changes to the domain. The maximum is 4000 characters.
* `email_notification_list` - (List) A list of email addresses to notify when a change is made to the domain.
* `default_timeout_penalty` - (Default: 25) Specifies the timeout penalty score.
* `load_imbalance_percentage` - Indicates the percent of load imbalance factor (LIF) for the domain.
* `default_ssl_client_private_key` - Specifies an optional Base64-encoded private key that corresponds with the TLS certificate for TLS-based liveness tests (HTTPS, SMTPS, POPS, and TCPS).
* `default_error_penalty` - (Default: 75) Specifies the download penalty score. If the download encounters an error, the web agent computes a score that is either the download time in seconds or a penalty score.
* `cname_coalescing_enabled` - (Boolean) If enabled, GTM collapses CNAME redirections in DNS answers when it knows the target of the CNAME.
* `load_feedback` - (Boolean) Indicates whether one or more measurements of load (resources) are defined by you and supplied by each data center in real time to balance load.
* `default_ssl_client_certificate` - Specifies an optional Base64-encoded certificate that corresponds with the private key for TLS-based liveness tests (HTTPS, SMTPS, POPS, and TCPS).
* `end_user_mapping_enabled` - (Boolean) Indicates whether the GTM Domain is using end user client subnet mapping.

### Computed

The following arguments will be found in `terraform.tfstate` and can be referenced throughout the configuration. The values cannot be changed.

* `default_unreachable_threshold`.
* `min_pingable_region_fraction`.
* `servermonitor_liveness_count`.
* `round_robin_prefix`.
* `servermonitor_load_count`.
* `ping_interval`.
* `max_ttl`.
* `default_health_max`.
* `map_update_interval`.
* `max_properties`.
* `max_resources`.
* `default_error_penalty`.
* `max_test_timeout`.
* `default_health_multiplier`.
* `servermonitor_pool`.
* `min_ttl`.
* `default_max_unreachable_penalty`.
* `default_health_threshold`.
* `min_test_interval`.
* `ping_packet_size`.

### Schema Reference

The GTM Domain backing schema and more complete element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#domain).
