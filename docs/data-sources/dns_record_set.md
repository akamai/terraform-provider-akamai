---
layout: akamai
subcategory: Edge DNS
---

# akamai_dns_record_set

Use the `akamai_dns_record_set` data source to list the record sets for a zone.

## Example usage

Basic usage:

```
data "akamai_dns_record_set" "test" {
  zone        = "exampleterraform.io"
  host        = "www.exampleterraform.io"
  record_type = "A"
}

output "test_addrs" {
  value = join(",", data.akamai_dns_record_set.test.rdata)
}
```

## Argument reference

This data source supports these arguments:

* `zone` - (Required) The domain zone, including any nested subdomains.
* `host` - (Required) The base credential hostname without the protocol.
* `record_type` - (Required) The DNS record type.

## Attributes reference

This data source supports this attribute:

* `rdata` - An array of data strings, representing multiple records within a set.
