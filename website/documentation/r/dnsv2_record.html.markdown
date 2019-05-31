---
layout: "akamai"
page_title: "Akamai: dnsv2 record"
sidebar_current: "docs-akamai-dnsv2_record"
description: |-
  DNS V2Record
---

# akamai_dnsv2_record


The `akamai_dnsv2_record` provides the resource for configuring a dns record to integrate easily with your existing DNS infrastructure to provide a secure, high performance, highly available and scalable solution for DNS hosting.



## Example Usage

Basic usage:

```hcl

#A record
resource "akamai_dnsv2_record" "origin" {

    zone = "akamaideveloper.net"
    name = "test.akamaideveloper.net"
    recordtype =  "A"
    active = true
    ttl =  30
    target = ["${aws_elb.demoapp.dns_name}"]

}
#CName record
resource "akamai_dnsv2_record" "edge" {

    zone = "example.com"
    name = "demo.example.com”
    recordtype =  "CNAME"
    active = true
    ttl =  600 
    target = “sample.edgesuite.com”
}



```

## Argument Reference

The following arguments are supported:

`*name` — (Required) The name of the record. The name is an owner name, that is, the name of the node to which this resource record pertains.  
`*zone` — (Required) Domain zone, encapsulating any nested subdomains.  
`*recordType` — (Required) The DNS record type.  
`*active` — (Required,Boolean) Whether the record is active.  
`*ttl` — (Required,Boolean) The TTL is a 32-bit signed integer that specifies the time interval that the resource record may be cached before the source of the information should be consulted again. Zero values are interpreted to mean that the RR can only be used for the transaction in progress, and should not be cached. Zero values can also be used for extremely volatile data.  
`*target` — (Required) A domain name that specifies the canonical or primary name for the owner. The owner name is an alias.  
