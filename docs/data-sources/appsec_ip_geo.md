---
layout: "akamai"
page_title: "Akamai: IP/Geo"
subcategory: "Application Security"
description: |-
 IP/Geo
---

# akamai_appsec_ip_geo

Use the `akamai_appsec_ip_geo` data source to retrieve information about which network lists are used in the IP/Geo Firewall settings.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to see IP/GEO firewall settings
data "akamai_appsec_ip_geo" "ip_geo" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}

output "ip_geo_mode" {
  value = akamai_appsec_ip_geo.ip_geo.mode 
}

output "geo_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo.geo_network_lists
}

output "ip_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo.ip_network_lists
}

output "exception_ip_network_lists" {
  value = akamai_appsec_ip_geo.ip_geo.exception_ip_network_lists
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Optional) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `mode` - The mode used for IP/Geo firewall blocking: `block` to block specific IPs, geographies or network lists, or `allow` to allow specific IPs or geographies to be let through while blocking the rest.

* `geo_network_lists` - The network lists to be blocked or allowed geographically.

* `ip_network_lists` - The network lists to be blocked or allowd by IP address.

* `exception_ip_network_lists` - The network lists to be allowed regardless of `mode`, `geo_network_lists`, and `ip_network_lists` parameters.

* `output_txt` - A tabular display of the IP/Geo firewall settings.

