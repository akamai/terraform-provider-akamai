---
layout: "akamai"
page_title: "Akamai: dnsv2 zone"
sidebar_current: "docs-akamai-dnsv2_zone"
description: |-
  DNS V2Zone
---

# akamai_dnsv2_zone


The `akamai_dnsv2_zone` provides the resource for configuring a dns zone to integrate easily with your existing DNS infrastructure to provide a secure, high performance, highly available and scalable solution for DNS hosting.



## Example Usage

Basic usage:

```hcl

resource "akamai_dnsv2_zone" "demozone" {

    contractid = "ctr_XXX"
    zone = "example.com"
    type =  "Primary”
    masters = [
              “1.2.3.4”,
              “1.2.3.5”
              ]
    comment =  “some comment”
    gid =100
    signandserve = true
}




```

## Argument Reference

The following arguments are supported:

*`contractid` — (Required) The contract ID.  
*`zone` — (Required) Domain zone, encapsulating any nested subdomains.  
*`type` — (Required) Whether the zone is primary or secondary.  
*`masters` — (Required) The names or addresses of the customer’s nameservers from which the zone data should be retrieved.  
*`comment` — (Required) The comment.  
*`gid` — (Required) The currently selected group ID.  
*`signandserve` — (Required) Whether DNSSEC Sign&Serve is enabled.  
