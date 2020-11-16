---
layout: "akamai"
page_title: "Akamai: edge hostname"
subcategory: "Provisioning"
description: |-
  Edge Hostname
---

# akamai_edge_hostname

The `akamai_edge_hostname` provides the resource for configuring a secure edge hostname that determines how requests for your site, app, or content are mapped to Akamai edge servers. 

An edge hostname is the CNAME target you use when directing your end user traffic to Akamai. In a typical DNS CNAME, your www.customer.com hostname corresponds to an edge hostname of www.customer.com.edgesuite.net.


## Example Usage

Basic usage:

```hcl
resource "akamai_edge_hostname" "terraform-demo" {
    product_id  = "prd_Object_Delivery"
    contract_id = "ctr_1-AB123"
    group_id    = "grp_123"
    edge_hostname = "www.example.org.edgesuite.net"
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The Edge Code name
* `contract_id` — (Required) The Contract ID.  Can be provided with or without `ctr_` prefix.
* `product_id` — (Required) The Product ID. Can be provided with or without `prd_` prefix.
* `edge_hostname` — (Required) One or more edge hostnames (must be <= to the number of public hostnames).
* `certificate` — (Optional) The certificate enrollment ID. Required when `edge_hostname` ends in edgekey.net.
* `ip_behavior` - (Optional) Directly specify IP protocol the hostname supports. Must be one of : `IPV4`, `IPV6_PERFORMANCE` or `IPV6_COMPLIANCE`

### Deprecated Arguments
* `contract` — (Deprecated) synonym of contract_id for legacy purposes
* `ipv4` — (Optional, Deprecated) Whether property supports IPv4 to origin. Used to compute `ip_behavior` when it is not supplied. (Default: `true`).
* `ipv6` —  (Optional, Deprecated) Whether property supports IPv6 to origin. Used to compute `ip_behavior` when it is not supplied. (Default: `false`).
* `product` — (Deprecated) synonym of product_id for legacy purposes

## Attributes Reference

The following attributes are returned:

* `ip_behavior` — One of `IPV4`, `IPV6` or `IPV6_COMPLIANCE` to specify whether hostname will use IPV4 IPV6 or both.* `ip_behavior` — One of `IPV4`, `IPV6_PERFORMANCE` or `IPV6_COMPLIANCE` to specify whether hostname will use IPV4 IPV6 or both.