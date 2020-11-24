---
layout: "akamai"
page_title: "Akamai: edge hostname"
subcategory: "Provisioning"
description: |-
  Edge Hostname
---

# akamai_edge_hostname

The `akamai_edge_hostname` resource lets you configure a secure edge hostname. Your edge hostname determines how requests for your site, app, or content are mapped to Akamai edge servers. 

An edge hostname is the CNAME target you use when directing your end user traffic to Akamai. Each hostname assigned to a property has a corresponding edge hostname. 
 
Akamai supports three types of edge hostnames, depending on the level of security you need for your traffic: Standard TLS, Enhanced TLS, and Shared Certificate. When entering the `edge_hostname` attribute, you need to include a specific domain suffix for your edge hostname type: 

| Edge hostname type | Domain suffix |
|------|-------|
| Enhanced TLS | edgekey.net |
| Standard TLS | edgesuite.net |
| Shared Cert | akamaized.net |

For example, if you use Standard TLS and have `www.customer.com` as a hostname, your edge hostname would be `www.customer.com.edgesuite.net`. If you wanted to use Enhanced TLS with the same hostname, your edge hostname would be `www.example.com.edgekey.net`. See the [Property Manager API (PAPI)](https://developer.akamai.com/api/core_features/property_manager/v1.html#createedgehostnames) for more information.


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

This resource supports these arguments:

* `name` - (Required) The name of the edge hostname.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 
* `product_id` - (Required) A product's unique ID, including the `prd_` prefix.
* `edge_hostname` - (Required) One or more edge hostnames. The number of edge hostnames must be less than or equal to the number of public hostnames.
* `certificate` - (Optional) Required only when creating an Enhanced TLS edge hostname. This argument sets the certificate enrollment ID. Edge hostnames (`edge_hostname`) for Enhanced TLS end in `edgekey.net`. You can retrieve this ID from the Certificate Provisioning System.
* `ip_behavior` - (Optional) Which version of the IP protocol to use: `IPV4` for version 4 only, `IPV6_PERFORMANCE` for version 6 only, or `IPV6_COMPLIANCE` for both 4 and 6. The default value is `IPV4`.

### Deprecated Arguments

* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `ipv6` -  (Deprecated) Optional argument used when a property supports IPv6 to origin. An existing resource will use this argument if `ip_behavior` hasn't been added. Set to `false` by default.
* `ipv4` - (Deprecated) Optional argument used when a property supports IPv4 to origin. An existing resource will use this argument if `ip_behavior` hasn't been added. Set to `true` by default.
* `product` - (Deprecated) Replaced by `product_id`. Maintained for legacy purposes.

## Attributes Reference

This resource returns this attribute:

* `ip_behavior` - Returns the IP protocol the hostname will use, either `IPV4` for version 4, IPV6_PERFORMANCE` for version 6, or `IPV6_COMPLIANCE` for both.