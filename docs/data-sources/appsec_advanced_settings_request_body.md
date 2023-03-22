---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_advanced_settings_request_body

**Scopes**: Security configuration; security policy

Returns information about your Request Size Inspection Limit controls. By default, information is returned for all the security policies in the configuration.
However, you can return data for a single policy by using the `security_policy_id` parameter.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/advanced-settings/request-body](https://techdocs.akamai.com/application-security/reference/get-policies-request-body)

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

// USE CASE: User wants to view the advanced settings within the Request Size Inspection Limit settings to access a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_advanced_settings_request_body" "request_body" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "advanced_settings_request_body_json" {
  value = data.akamai_appsec_advanced_settings_request_body.request_body.json
}

output "advanced_settings_request_body_output" {
  value = data.akamai_appsec_advanced_settings_request_body.request_body.output_text
}

data "akamai_appsec_advanced_settings_request_body" "policy_override" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "advanced_settings_policy_request_body_output" {
  value = data.akamai_appsec_advanced_settings_request_body.policy_override.output_text
}

output "advanced_settings_policy_request_body_json" {
  value = data.akamai_appsec_advanced_settings_request_body.policy_override.json
}

```
## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the Request Size Inspection Limit settings.
- `security_policy_id` (Optional). Unique identifier of the security policy associated with the Request Size Inspection Limit settings. If not included, information is returned for all your security policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of information about the Request Size Inspection Limit settings.
- `output_text`. Tabular report showing the Request Size Inspection Limit settings.