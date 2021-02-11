---
layout: "akamai"
page_title: "Akamai: gtm domain"
subcategory: "Global Traffic Management"
description: |-
  GTM Domain
---

# akamai_gtm_domain

Use the `akamai_gtm_domain` resource to create, configure, and import a GTM Domain, which is a basic building block of a traffic management configuration. 

~> **Note** Import requires an ID with this format: `existing_domain_name`.

## Example usage

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

## Argument reference

This resource supports these arguments:

* `contract` - (Required) If creating a domain, the contract ID.
* `group` - (Required) If creating a domain, the currently selected group ID.
* `name` - (Required) The DNS name for a collection of GTM Properties.
* `type` - (Required) Th type of GTM domain. Options include `failover-only`, `static`, `weighted`, `basic`, or `full`. 
* `wait_on_complete` - (Optional) A boolean that, if set to `true`, waits for transaction to complete.
* `comment` - (Optional) A descriptive note about changes to the domain. The maximum is 4000 characters.
* `email_notification_list` - (Optional) A list of email addresses to notify when a change is made to the domain.
* `default_timeout_penalty` - (Optional) Specifies the timeout penalty score. Default is `25`.
* `load_imbalance_percentage` - (Optional) Indicates the percentage of load imbalance factor (LIF) for the domain.
* `default_ssl_client_private_key` - (Optional) Specifies a Base64-encoded private key that corresponds with the TLS certificate for HTTPS, SMTPS, POPS, and TCPS liveness tests.
* `default_error_penalty` - (Optional) Specifies the download penalty score. The default is `75`. If the download encounters an error, the web agent computes a score that is either the download time in seconds or a penalty score.
* `cname_coalescing_enabled` - (Optional) A boolean that if set to `true`, GTM collapses CNAME redirections in DNS answers when it knows the target of the CNAME.
* `load_feedback` - (Optional) A boolean indicating whether one or more measurements of load (resources) are defined by you and supplied by each data center in real time to balance load.
* `default_ssl_client_certificate` - (Optional) Specifies an optional Base64-encoded certificate that corresponds with the private key for TLS-based liveness tests (HTTPS, SMTPS, POPS, and TCPS).
* `end_user_mapping_enabled` - (Optional) A boolean indicating whether whether the GTM Domain is using end user client subnet mapping.

### Computed argument reference

This resource returns these computed arguments in the `terraform.tfstate` file:

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

### Schema reference

You can download the GTM Domain backing schema from the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#domain) page.
