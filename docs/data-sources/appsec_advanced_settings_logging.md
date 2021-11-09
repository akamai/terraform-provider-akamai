---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsLogging"
subcategory: "Application Security"
description: |-
 AdvancedSettingsLogging
---

# akamai_appsec_advanced_settings_logging

**Scopes**: Security configuration; security policy

Returns information about your HTTP header logging controls. By default, information is returned for all the security policies in the configuration; however, you can return data for a single policy by using the `security_policy_id` parameter. The returned information is described in the [ConfigHeaderLog members](https://developer.akamai.com/api/cloud_security/application_security/v1.html#a6d1c316) section of the Application Security API.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/logging](https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethttpheaderloggingforaconfiguration)

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

// USE CASE: User wants to view the custom rules associated with a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_custom_rules" "custom_rules" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
output "custom_rules_output_text" {
  value = data.akamai_appsec_custom_rules.custom_rules.output_text
}
output "custom_rules_json" {
  value = data.akamai_appsec_custom_rules.custom_rules.json
}
output "custom_rules_config_id" {
  value = data.akamai_appsec_custom_rules.custom_rules.config_id
}
// USE CASE: User wants to view a specific custom rule.

data "akamai_appsec_custom_rules" "specific_custom_rule" {
  config_id      = data.akamai_appsec_configuration.configuration.config_id
  custom_rule_id = "60029316"
}
output "specific_custom_rule_json" {
  value = data.akamai_appsec_custom_rules.specific_custom_rule.json
}
```
## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the logging settings.
- `security_policy_id` (Optional). Unique identifier of the security policy associated with the logging settings. If not included, information is returned for all your security policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of information about the logging settings.
- `output_text`. Tabular report showing the logging settings.

