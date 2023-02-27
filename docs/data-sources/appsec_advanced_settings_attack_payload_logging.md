---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_advanced_settings_attack_payload_logging

**Scopes**: Security configuration; security policy

Returns information about your Attack Payload Logging controls. By default, information is returned for all the security policies in the configuration.
However, you can return data for a single policy by using the `security_policy_id` parameter.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/logging/attack-payload]

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

// USE CASE: User wants to view the advanced settings within the Attack Payload Logging settings to access a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_advanced_settings_attack_payload_logging" "attack_payload_logging" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "advanced_settings_attack_payload_logging_json" {
  value = data.akamai_appsec_advanced_settings_attack_payload_logging.attack_payload_logging.json
}

output "advanced_settings_attack_payload_logging_output" {
  value = data.akamai_appsec_advanced_settings_attack_payload_logging.attack_payload_logging.output_text
}

data "akamai_appsec_advanced_settings_attack_payload_logging" "policy_override" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "advanced_settings_policy_attack_payload_logging_output" {
  value = data.akamai_appsec_advanced_settings_attack_payload_logging.policy_override.output_text
}

output "advanced_settings_policy_attack_payload_logging_json" {
  value = data.akamai_appsec_advanced_settings_attack_payload_logging.policy_override.json
}

```
## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the Attack Payload Logging settings.
- `security_policy_id` (Optional). Unique identifier of the security policy associated with the Attack Payload Logging settings. If not included, information is returned for all your security policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of information about the Attack Payload Logging settings.
- `output_text`. Tabular report showing the Attack Payload Logging settings.