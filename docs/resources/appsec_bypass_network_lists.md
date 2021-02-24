---
layout: "akamai"
page_title: "Akamai: BypassNetworkLists"
subcategory: "Application Security"
description: |-
 BypassNetworkLists
---

# akamai_appsec_bypass_network_lists

Use the `akamai_appsec_bypass_network_lists` resource to update which network lists to use in the bypass network lists settings.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to update bypass network lists used in a Security Configuration version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_bypass_network_lists" "bypass_network_lists" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  bypass_network_list = ["id1","id2"]
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `version` - (Required) The version number of the configuration to use.

* `bypass_network_list` - (Required) A list containing the IDs of the network lists to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the updated list of network list IDs.

