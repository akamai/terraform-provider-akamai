---
layout: "akamai"
page_title: "Akamai: SelectableHostnames"
subcategory: "Application Security"
description: |-
 SelectableHostnames
---

# akamai_appsec_selectable_hostnames

Use the `akamai_appsec_selectable_hostnames` data source to retrieve the list of hostnames that may be protected under a given security configuration version. You can use specify the list to be retrieved either by supplying the name and version of a security configuration, or by supplying a group ID and contract ID.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the list of hosts available to be added to the list of those protected
//           under a given security configuration and version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

output "selectable_hostnames" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

// USE CASE: user wants to view the same list of unprotected hostnames, in JSON form
output "selectable_hostnames_json" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames_json
}

// USE CASE: user wants to view the same list of unprotected hostnames, in tabular form
output "selectable_hostnames_output_text" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.output_text
}

//USE CASE: user wants to view the list of hosts available to create a new config under a given contractid and groupid
data "akamai_appsec_selectable_hostnames" "selectable_hostnames_for_create_configuration" {
  contractid = var.contractid
  groupid = var.groupid
}

output "selectable_hostnames_for_create_configuration" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.hostnames
}

// USE CASE: user wants to view the same list of available hostnames, in JSON form
output "selectable_hostnames_for_create_configuration_json" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.hostnames_json
}

// USE CASE: user wants to view the same list of available hostnames, in tabular form
output "selectable_hostnames_for_create_configuration_output_text" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Optional) The ID of the security configuration to use.

* `version` - (Optional) The version number of the security configuration to use.

* `contractid` - (Optional) The ID of the contract to use.

* `groupid` - (Optional) The ID of the group to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `hostnames` - The list of selectable hostnames.

* `hostnames_json` - The list of selectable hostnames in json format.

* `output_text` - A tabular display of the selectable hostnames showing the name and config_id of the security configuration under which the host is protected in production, or '-' if the host is not protected in production.

