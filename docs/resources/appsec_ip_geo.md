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

//OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#putipgeofirewall

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

// USE CASE: user wants to update the IP/GEO firewall mode to "block all traffic except IPs/Subnets in block execptions" and update the Exception list
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

//TF destroy - should call protections protections API and turn it off the ipGeo control off
//OPEN API: /appsec/v1/configs/{config_id}/versions/{version}/security-policies/{security_policy_id}/protections
//Request body: {"applyNetworkLayerControls":false}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `mode` - (Required) The mode to use for IP/Geo firewall blocking: `blockSpecificIPGeo` to block specific IPs, geographies or network lists, or `blockAllTrafficExceptAllowedIPs` to allow specific IPs or geographies to be let through while blocking the rest.

* `geo_network_lists` - (Optional) TBD

* `ip_network_lists` - (Optional) TBD

* `exception_ip_network_lists` - (Required) TBD:w

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None
