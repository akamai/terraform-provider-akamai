---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsLogging"
subcategory: "Application Security"
description: |-
 AdvancedSettingsLogging
---

# akamai_appsec_advanced_settings_logging

Use the `akamai_appsec_advanced_settings_logging` data source to retrieve information about the HTTP header logging controls for a configuration. This operation applies at the configuration level, and therefore applies to all policies within a configuration. You may retrieve these settings for a particular policy by specifying the policy using the security_policy_id parameter. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethttpheaderloggingforaconfiguration).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the advanced settings loggings in a given security configuration
// when policy is set -  /appsec/v1/configs/{configId}/versions/{versionNum}/security-policies/{policyId}/advanced-settings/logging
// with out policy - /appsec/v1/configs/{configId}/versions/{versionNum}/advanced-settings/logging
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_advanced_settings_logging" "logging" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "advanced_settings_logging_output" {
  value = data.akamai_appsec_advanced_settings_logging.logging.output_text
}

output "advanced_settings_logging_json" {
  value = data.akamai_appsec_advanced_settings_logging.logging.json
}

data "akamai_appsec_advanced_settings_logging" "policy_override" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "advanced_settings_policy_logging_output" {
  value = data.akamai_appsec_advanced_settings_logging.policy_override.output_text
}

output "advanced_settings_policy_logging_json" {
  value = data.akamai_appsec_advanced_settings_logging.policy_override.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID.

* `security_policy_id` - (Optional) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of information about the logging settings.

* `output_text` - A tabular display showing the logging settings.

