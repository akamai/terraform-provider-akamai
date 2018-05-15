---
layout: "akamai"
page_title: "Akamai: fastdns_zone"
sidebar_current: "docs-akamai-resource-fastdns-zone"
description: |-
  Create and update Akamai Fast DNS Zone records
---

# akamai_fastdns_zone

The `akamai_fastdns_zone` resource represents an Akamai Fast DNS zone configuration, allowing you to update, and create records within that Zone.

## Example Usage

Basic usage:

```hcl
resource "akamai_fastdns_zone" "example_net" {
  hostname = "example.net"

  a {
    name = "www-origin"
    ttl = 600
    active = true
    target = "192.0.2.3"
  }
  
  cname {
    name = "www"
    ttl = 600
    active = true
    target = "www.example.net.edgesuite.net."
  }
}
```

## Argument Reference

The following arguments are supported:

* `hostname` — (Required) the zone hostname
* Common arguments for all record types:
  * `active` — (Optional, boolean, except SOA) Whether the record is active
  * `name` — (Required, except SOA) The record name
  * `ttl` — (Optional) Time To Live (in seconds)
* `a` — (Optional) An A record
  * `target` — (Required) Record target Record target
* `aaaa` — (Optional) An AAAA record
  * `target` — (Required) Record target
* `afsdb` — (Optional) An AFSDB record
  * `subtype` — (Required)
  * `target` — (Required) Record target
* `cname` — (Optional) An CNAME record
  * `target` — (Required) Record target
* `dnskey` — (Optional) An DNSKEY record
  * `algorithm` — (Required)
  * `flags` — (Required)
  * `key` — (Required)
  * `protocol` — (Required)
* `ds` — (Optional) An DS record
  * `algorithm` — (Required)
  * `digest` — (Required)
  * `digest-type` — (Required)
  * `keytag` — (Required)
* `hinfo` — (Optional) An HINFO record
  * `hardware` — (Required)
  * `software` — (Required)
* `loc` — (Optional) An LOC record
  * `target` — (Required) Record target
* `mx` — (Optional) An MX record
  * `priority` — (Required)
  * `target` — (Required) Record target
* `naptr` — (Optional) An NAPTR record
  * `flags` — (Required)
  * `order` — (Required)
  * `preference` — (Required)
  * `regexp` — (Required)
  * `replacement` — (Required)
  * `service` — (Required)
* `ns` — (Optional) An NS record
  * `target` — (Required) Record target
* `nsec3` — (Optional) An NSEC3 record
  * `algorithm` — (Required)
  * `flags` — (Required)
  * `iterations` — (Required)
  * `next-hashed-owner-name` — (Required)
  * `salt` — (Required)
  * `type-bitmaps` — (Required)
* `nsec3param` — (Optional) An NSEC3PARAM record
  * `algorithm` — (Required)
  * `flags` — (Required)
  * `iterations` — (Required)
  * `salt` — (Required)
* `ptr` — (Optional) An PTR record
  * `target` — (Required) Record target
* `rp` — (Optional) An RP record
  * `mailbox` — (Required)
  * `txt` — (Required)
* `rrsig` — (Optional) An RRSIG record
  * `algorithm` — (Required)
  * `expiration` — (Required)
  * `inception` — (Required)
  * `keytag` — (Required)
  * `labels` — (Required)
  * `original-ttl` — (Required)
  * `signature` — (Required)
  * `signer` — (Required)
  * `type-covered` — (Required)
* `soa` — (Optional) An SOA record
  * `contact` — (Required)
  * `expire` — (Required)
  * `minimum` — (Required)
  * `originserver` — (Required)
  * `refresh` — (Required)
  * `retry` — (Required)
  * `serial` — (Optional)
* `spf` — (Optional) An SPF record
  * `target` — (Required) Record target
* `srv` — (Optional) An SRV record
  * `port` — (Required)
  * `priority` — (Required)
  * `target` — (Required) Record target
  * `weight` — (Required)
* `sshfp` — (Optional) An SSHFP record
  * `algorithm` — (Required)
  * `fingerprint` — (Required)
  * `fingerprint-type` — (Required)
* `txt` — (Optional) An TXT record
  * `target` — (Required) Record target
  
A full description of each option is available [here](https://developer.akamai.com/api/luna/config-dns/data.html).