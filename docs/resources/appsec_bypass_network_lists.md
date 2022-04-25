---
layout: "akamai"
page_title: "Akamai: BypassNetworkLists"
subcategory: "Application Security"
description: |-
 BypassNetworkLists
---

# akamai_appsec_bypass_network_lists

**Scopes**: Security configuration

Specifies the networks that appear on the bypass network list. Networks on this list are allowed to bypass the Web Application Firewall.

Note that this resource is only applicable to WAP (Web Application Protector) configurations.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/bypass-network-lists](https://techdocs.akamai.com/application-security/reference/put-bypass-network-lists)

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

// USE CASE: User wants to update the bypass network list used in a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_bypass_network_lists" "bypass_network_lists" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  bypass_network_list = ["DocumentationNetworkList", "TrainingNetworkList"]
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bypass network lists being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bypass network lists.
- `bypass_network_list` (Required). JSON array of network IDs that comprise the bypass list.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the updated list of bypass network IDs.