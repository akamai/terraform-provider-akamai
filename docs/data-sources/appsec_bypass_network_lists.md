---
layout: "akamai"
page_title: "Akamai: BypassNetworkLists"
subcategory: "Application Security"
description: |-
 BypassNetworkLists
---

# akamai_appsec_bypass_network_lists

Use the `akamai_appsec_bypass_network_lists` data source to retrieve information about which network
lists are used in the bypass network lists settings.  The information available is described
[here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getbypassnetworklistsforawapconfigversion).
Note: this data source is only applicable to WAP (Web Application Protector) configurations.


## Example Usage

Basic usage:

```hcl

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to see information about bypass network lists used in a Security Configuration version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_bypass_network_lists" "bypass_network_lists" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

//Tabular display of ID and Name of the network lists 
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

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `bypass_network_list` - A list of strings containing the network list IDs.

* `json` - A JSON-formatted list of information about the bypass network lists.

* `output_text` - A tabular display showing the bypass network list information.

