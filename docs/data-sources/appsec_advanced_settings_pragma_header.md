---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPragmaHeader"
subcategory: "Application Security"
description: |-
 AdvancedSettingsPragmaHeader
---

# akamai_appsec_advanced_settings_pragma_header

Use the `akamai_appsec_advanced_settings_pragma_header` data source to retrieve pragma header settings for a configuration or a security policy. Additional information is available [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpragmaheaderconfiguration).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the advanced settings for pragma header in a given security configuration or a security policy
// when policy is set -  /appsec/v1/configs/{configId}/versions/{versionNum}/security-policies/{policyId}/advanced-settings/pragma-header
// without policy - /appsec/v1/configs/{configId}/versions/{versionNum}/advanced-settings/pragma-header
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
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
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "advanced_settings_policy_pragma_header_output" {
  value = data.akamai_appsec_advanced_settings_pragma_header.policy_pragma_header.output_text
}

output "advanced_settings_policy_pragma_header_json" {
  value = data.akamai_appsec_advanced_settings_pragma_header.policy_pragma_header.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID.

* `security_policy_id` - (Optional) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpragmaheaderpolicy)) list of information about the pragma header settings.

* `output_text` - A tabular display showing the pragma header settings.

