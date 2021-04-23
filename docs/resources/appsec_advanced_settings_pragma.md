---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPragmaHeader"
subcategory: "Application Security"
description: |-
  AdvancedSettingsPragmaHeader
---

# resource_akamai_appsec_advanced_settings_pragma_header

The `resource_akamai_appsec_advanced_settings_pragma_header` resource allows you to specify which headers you can exclude from inspection when you pass a `Pragma` debug header. This operation applies at the configuration level or policy level.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the pragma header settings for a configuration or a security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_advanced_settings_pragma_header" "pragma_header" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  pragma_header = file("${path.module}/pragma_header.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Optional) The ID of the security policy to use.

* `pragma_header` - (Required) The name of a file containing a JSON-formatted ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpragmaheaderpolicy)) description of the conditions to exclude from the default `remove` action. By default, the Pragma header debugging information is stripped from an operation’s response except in cases where you set excludeCondition. To remove existing settings, submit your request with an empty payload {} at the top-level of an object. For example, submit "type": "{}" in the request body to remove the REQUEST_HEADER_VALUE_MATCH from the excluded conditions. If you submit an empty payload for each member, you’ll clear all of your condition settings. To modify Pragma header settings at the security configuration level, run Modify Pragma header settings for a configuration. Contact your account team if you’d like to run this operation.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None
