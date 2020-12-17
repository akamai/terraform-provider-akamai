---
layout: "akamai"
page_title: "Akamai: IP/Geo Firewall"
subcategory: "Application Security"
description: |-
 IP/Geo Firewall
---

# akamai_appsec_ip_geo

Use the `akamai_appsec_ip_geo` resource to update the method and which network lists to use for IP/Geo firewall blocking.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to update the IP/GEO firewall mode to "block specific IPs/Subnets and Geos" and update the IP list, GEO list & Exception list
resource  "akamai_appsec_ip_geo" "ip_geo_block" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id1
  mode = var.block
  geo_network_lists= var.geo_network_lists
  ip_network_lists= var.ip_network_lists
  exception_ip_network_lists= var.exception_ip_network_lists
}

// USE CASE: user wants to update the IP/GEO firewall mode to "block all traffic except IPs/Subnets in block exceptions" and update the Exception list
resource  "akamai_appsec_ip_geo" "ip_geo_allow" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id2
  mode = var.allow
  exception_ip_network_lists= var.exception_ip_network_lists
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

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `mode` - (Required) The mode to use for IP/Geo firewall blocking: `block` to block specific IPs, geographies or network lists, or `allow` to allow specific IPs or geographies to be let through while blocking the rest.

* `geo_network_lists` - (Optional) The network lists to be blocked or allowed geographically.

* `ip_network_lists` - (Optional) The network lists to be blocked or allowd by IP address.

* `exception_ip_network_lists` - (Required) The network lists to be allowed regardless of `mode`, `geo_network_lists`, and `ip_network_lists` parameters.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_txt` - A tabular display of the IP/Geo firewall settings.
