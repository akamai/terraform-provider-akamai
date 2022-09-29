---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_ip_geo

**Scopes**: Security policy

Modifies the method used for firewall blocking, and manages the network lists used for IP/Geo firewall blocking.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/ip-geo-firewall](https://techdocs.akamai.com/application-security/reference/put-policy-ip-geo-firewall)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

// USE CASE: User wants to update the IP/Geo firewall mode, and update the IP, geographic, and exception lists.

resource "akamai_appsec_ip_geo" "ip_geo_block" {
  config_id                  = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id         = "gms1_134637"
  mode                       = "block"
  geo_network_lists          = ["06038_GEO_TEST"]
  ip_network_lists           = ["56921_TEST"]
  exception_ip_network_lists = ["07126_EXCEPTION_TEST"]
}

// USE CASE: User wants to update the IP/Geo firewall mode and update the exception list.

resource "akamai_appsec_ip_geo" "ip_geo_allow" {
  config_id                  = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id         = "gms1-090334"
  mode                       = "allow"
  exception_ip_network_lists = ["07126_EXCEPTION_TEST"]
}

output "ip_geo_mode_block" {
  value = akamai_appsec_ip_geo.ip_geo_block.mode
}

output "block_geo_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo_block.geo_network_lists
}

output "block_ip_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo_block.ip_network_lists
}

output "block_exception_ip_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo_block.exception_ip_network_lists
}

output "ip_geo_mode_allow" {
  value = akamai_appsec_ip_geo.ip_geo_allow.mode
}
output "allow_exception_ip_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo_allow.exception_ip_network_lists
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the IP/Geo lists being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the IP/Geo lists being modified.
- `mode` (Required). Set to **block** to prevent the specified network lists from being allowed through the firewall: all other entities will be allowed to pass through the firewall. Set to **allow** to allow the specified network lists to pass through the firewall; all other entities will be prevented from passing through the firewall.
- `geo_network_lists` (Optional). JSON array of geographic network lists that, depending on the value of the `mode` argument, will be blocked or allowed through the firewall.
- `ip_network_lists` (Optional). JSON array of IP network lists that, depending on the value of the `mode` argument, will be blocked or allowed through the firewall.
- `exception_ip_network_lists` (Optional). JSON array of network lists that are always allowed to pass through the firewall, regardless of the value of any other setting.