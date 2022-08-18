---
layout: "akamai"
page_title: "Akamai: CustomDeny"
subcategory: "Application Security"
description: |-
 CustomDeny
---


# akamai_appsec_custom_deny

**Scopes**: Security configuration; custom deny

Returns information about custom deny actions. Custom denies allow you to craft your own error messages or redirect pages to use when HTTP requests are denied.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-deny](https://techdocs.akamai.com/application-security/reference/get-custom-deny-actions)

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

// USE CASE: User wants to view the custom deny data for a given security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_custom_deny" "custom_deny_list" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

// USE CASE: User wants to display the returned data in a table.

output "custom_deny_list_output" {
  value = data.akamai_appsec_custom_deny.custom_deny_list.output_text
}

output "custom_deny_list_json" {
  value = data.akamai_appsec_custom_deny.custom_deny_list.json
}

// USE CASE: User wants to view a specific custom deny associated with a security configuration.

data "akamai_appsec_custom_deny" "custom_deny" {
  config_id      = data.akamai_appsec_configuration.configuration.config_id
  custom_deny_id = "deny_custom_64386"
}

output "custom_deny_json" {
  value = data.akamai_appsec_custom_deny.custom_deny.json
}

output "custom_deny_output" {
  value = data.akamai_appsec_custom_deny.custom_deny.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom denies.
- `custom_deny_id` (Optional). Unique identifier of the custom deny you want to return information for. If not included. information is returned for all your custom denies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of custom deny information.
- `output_text`. Tabular report of the custom deny information.