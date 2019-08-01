---
layout: "akamai"
page_title: "Akamai: edge hostname"
sidebar_current: "docs-akamai-resource-edge-hostname"
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
    product  = "prd_####"
    contract = "ctr_####"
    group    = "grp_####"
    edge_hostname = "www.example.org.edgesuite.net"
}
```

## Argument Reference

The following arguments are supported:

* `contract` — (Required) The contract ID.  
* `group` — (Required) The group ID.  
* `product` — (Required) The product ID.  
* `edge_hostname` — (Required) One or more edge hostnames (must be <= to the number of public hostnames).
* `ipv4` - (Optional) Whether the property supports IPv4 to origin.  (Default: `true`).
* `ipv6` —  (Optional) Whether the property supports IPv6 to origin. (Default: `false`).
* `certificate` - (Optional) The certificate enrollment ID.  

## Attributes Reference

The following attributes are returned:

* `ip_behavior`: Whether the hostname uses `IPV4`, `IPV6` or `IPV6_COMPLIANCE`.