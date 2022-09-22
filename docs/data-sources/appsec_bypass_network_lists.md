---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_bypass_network_lists

**Scopes**: Security configuration

Returns information about the network lists assigned to the bypass network list; networks on this list are not subject to firewall checking. 

Note that this data source is only applicable to WAP (Web Application Protector) configurations.

**Related API Endpoint**:[/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/bypass-network-lists](https://techdocs.akamai.com/application-security/reference/get-bypass-network-lists-per-policy)

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

// USE CASE: User wants to view information about the bypass network list used in a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_bypass_network_lists" "bypass_network_lists" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

// USE CASE: User wants to display returned data in a table.

output "bypass_network_lists_output" {
  value = data.akamai_appsec_bypass_network_lists.bypass_network_lists.output_text
}

output "bypass_network_lists_json" {
  value = data.akamai_appsec_bypass_network_lists.bypass_network_lists.json
}

output "bypass_network_lists_id_list" {
  value = data.akamai_appsec_bypass_network_lists.bypass_network_lists.bypass_network_list
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bypass network lists.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bypass network lists.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `bypass_network_list`. List of network IDs.
- `json`. JSON-formatted list of information about the bypass networks.
- `output_text`. Tabular report showing the bypass network list information.