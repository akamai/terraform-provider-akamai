---
layout: "akamai"
page_title: "Akamai: dns zone"
subcategory: "DNS"
description: |-
  DNS Zone
---

# akamai_dns_zone

The `akamai_dns_zone` provides the resource for configuring a dns zone to integrate easily with your existing DNS infrastructure to provide a secure, high performance, highly available and scalable solution for DNS hosting.

## Example Usage

Basic usage:

```hcl
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

## Argument Reference

The following arguments are supported:

* `contract` — (Required) The contract ID. 
* `group` — (Required) The currently selected group ID.   
* `zone` — (Required) Domain zone, encapsulating any nested subdomains.  
* `type` — (Required) Whether the zone is primary or secondary.  
* `masters` — (Required for Secondary) The names or addresses of the customer’s nameservers from which the zone data should be retrieved.  
* `comment` — (Required) A descriptive comment.  
* `sign_and_serve` — (Optional) Whether DNSSEC Sign&Serve is enabled. 
* `sign_and_serve_algorithm` — (Optional) Algorithm used by Sign&Serve.
* `target` — (Required for Alias) 
* `tsig_key` — (Optional) TSIG Key used in secure zone transfers
  * `name` - key name
  * `algorithm`
  * `secret`
* `end_customer_id` — (Optional)
  
