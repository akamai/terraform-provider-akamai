---
layout: "akamai"
page_title: "Akamai: IP/Geo"
subcategory: "Application Security"
description: |-
 IP/Geo
---


# akamai_appsec_ip_geo

**Scopes**: Security configuration; security policy

Returns information about the network lists used in the IP/Geo Firewall settings; also returns the firewall `mode`, which indicates whether devices on the geographic or IP address lists are allowed through the firewall or are blocked by the firewall.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/ip-geo-firewall](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getipgeofirewall)

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

// USE CASE: User wants to view IP/Geo firewall settings.

data "akamai_appsec_ip_geo" "ip_geo" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "ip_geo_mode" {
  value = data.akamai_appsec_ip_geo.ip_geo.mode
}

output "geo_network_lists" {
  value = data.akamai_appsec_ip_geo.ip_geo.geo_network_lists
}

output "ip_network_lists" {
  value = data.akamai_appsec_ip_geo.ip_geo.ip_network_lists
}

output "exception_ip_network_lists" {
  value = data.akamai_appsec_ip_geo.ip_geo.exception_ip_network_lists
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the IP/Geo lists.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the IP/Geo lists. If not included, information is returned for all your security policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `mode`. Specifies the action taken by the IP/Geo firewall. Valid values are:
  - **block**. Networks on the IP and geographic network lists are prevented from passing through the firewall.
  - **allow**.  Networks on the IP and geographic network lists are allowed to pass through the firewall.
- `geo_network_lists`. Network lists blocked or allowed based on geographic location.
- `ip_network_lists`. Network lists blocked or allowed based on IP address.
- `exception_ip_network_lists`. Network lists allowed through the firewall regardless of the values assigned to the `mode`, `geo_network_lists`, and `ip_network_lists` parameters.
- `output_text`. Tabular report of the IP/Geo firewall settings.

