---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPragmaHeader"
subcategory: "Application Security"
description: |-
 AdvancedSettingsPragmaHeader
---

# akamai_appsec_advanced_settings_pragma_header

**Scopes**: Security configuration; security policy

Returns pragma header settings information. This HTTP header provides information about such things as: the edge routers used in a transaction; the Akamai IP addresses involved; information about whether a request was cached or not; and so on. By default, pragma headers are removed from all responses.

Additional information is available from the [PragmaHeader members](https://developer.akamai.com/api/cloud_security/application_security/v1.html#64c92ba1) section of the Application Security API.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/pragma-header](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpragmaheaderconfiguration)

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

// USE CASE: User wants to view the pragma header settings for a security configuration or security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_advanced_settings_pragma_header" "pragma_header" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "advanced_settings_pragma_header_output" {
  value = data.akamai_appsec_advanced_settings_pragma_header.pragma_header.output_text
}

output "advanced_settings_pragma_header_json" {
  value = data.akamai_appsec_advanced_settings_pragma_header.pragma_header.json
}

data "akamai_appsec_advanced_settings_pragma_header" "policy_pragma_header" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "advanced_settings_policy_pragma_header_output" {
  value = data.akamai_appsec_advanced_settings_pragma_header.policy_pragma_header.output_text
}

output "advanced_settings_policy_pragma_header_json" {
  value = data.akamai_appsec_advanced_settings_pragma_header.policy_pragma_header.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the pragma header settings.
- `security_policy_id` (Optional). Unique identifier of the security policy associated with the pragma header settings. If not included, information is returned for all your security policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of information about the pragma header settings.
- `output_text`. Tabular report showing the pragma header settings.

