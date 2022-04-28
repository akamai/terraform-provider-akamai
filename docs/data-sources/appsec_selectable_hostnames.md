---
layout: "akamai"
page_title: "Akamai: SelectableHostnames"
subcategory: "Application Security"
description: |-
 SelectableHostnames
---

# akamai_appsec_selectable_hostnames

**Scopes**: Security configuration; contract; group

Returns the list of hostnames that can be (but aren't yet) protected by a security configuration. You can specify the set of hostnames to be retrieved either by supplying the name of a security configuration or by supplying an Akamai group ID and contract ID.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/selectable-hostnames](https://techdocs.akamai.com/application-security/reference/get-selectable-hostnames)

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

// USE CASE: User wants to view the hosts that can be protected by a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "selectable_hostnames" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

// USE CASE: User wants to view all the unprotected hostnames.

output "selectable_hostnames_json" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames_json
}

// USE CASE: user wants to view the same list of unprotected hostnames, in tabular form

output "selectable_hostnames_output_text" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.output_text
}

//USE CASE: User wants to view the list of hosts available for the specified contract and contract group before creating a new security configuration.

data "akamai_appsec_selectable_hostnames" "selectable_hostnames_for_create_configuration" {
  contract_id = "5-2WA382"
  group_id    = 12198
}

output "selectable_hostnames_for_create_configuration" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.hostnames
}

// USE CASE: User wants to view the available hostnames in JSON format.

output "selectable_hostnames_for_create_configuration_json" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.hostnames_json
}

// USE CASE: User wants to view the available hostnames in a table.

output "selectable_hostnames_for_create_configuration_output_text" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Optional). Unique identifier of the security configuration you want to return hostname information for. If not included, information is returned for all your security configurations. Note that argument can't be used with either the `contractid` or the `groupid` arguments.
- `contractid` (Optional). Unique identifier of the Akamai contract you want to return hostname information for. If not included, information is returned for all the Akamai contracts associated with your account. Note that this argument can't be used with the `config_id` argument.
- `groupid` (Optional). Unique identifier of the contract group you want to return hostname information for. If not included, information is returned for all your contract groups. (Or, if you include the `contractid` argument, all the groups associated with the specified contract.) Note that this argument can't be used with the `config_id` argument.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `hostnames`. List of selectable hostnames.
- `hostnames_json`. JSON-formatted list of selectable hostnames.
- `output_text`. Tabular report of the selectable hostnames showing the name and config_id of the security configuration under which the host is protected in production.