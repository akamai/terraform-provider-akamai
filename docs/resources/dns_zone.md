---
layout: "akamai"
page_title: "Akamai: dns zone"
subcategory: "DNS"
description: |-
  DNS Zone
---

# akamai_dns_zone

Use the `akamai_dns_zone` resource to configure a DNS zone that integrates with your existing DNS infrastructure.

## Example usage

Basic usage:

```
resource "akamai_dns_zone" "demozone" {
    contract = "ctr_1-AB123"
    group = 100
    zone = "example.com"
    type =  "secondary"
    masters = [
      "1.2.3.4",
      "1.2.3.5"
    ] 
    comment =  "some comment"
    sign_and_serve = false
}
```

## Argument reference

This resource supports these arguments:

* `comment` - (Required) A descriptive comment.
* `contract` - (Required) The contract ID.
* `group` - (Required) The currently selected group ID.
* `zone` - (Required) The domain zone, encapsulating any nested subdomains.
* `type` - (Required) Whether the zone is `primary`, `secondary`, or `alias`.
* `masters` - (Required for `secondary` zones) The names or IP addresses of the nameservers that the zone data should be retrieved from.
* `target` - (Required for `alias` zones) The name of the zone whose configuration this zone will copy.
* `sign_and_serve` - (Optional) Whether DNSSEC Sign and Serve is enabled.
* `sign_and_serve_algorithm` - (Optional) The algorithm used by Sign and Serve.
* `tsig_key` - (Optional) The TSIG Key used in secure zone transfers. If used, requires these arguments:
    * `name` - The key name.
    * `algorithm` - The hashing algorithm.
    * `secret` - String known between transfer endpoints.
* `end_customer_id` - (Optional) A free form identifier for the zone.
