---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsLogging"
subcategory: "Application Security"
description: |-
  AdvancedSettingsLogging
---

# resource_akamai_appsec_advanced_settings_logging

The `resource_akamai_appsec_advanced_settings_logging` resource allows you to enable, disable, or update HTTP header logging settings for a configuration. This operation applies at the configuration level, and therefore applies to all policies within a configuration. You may override these settings for a particular policy by specifying the policy using the security_policy_id parameter.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the logging settings
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_advanced_settings_logging" "logging" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  logging = file("${path.module}/logging.json")
}

// USE CASE: user wants to override the logging settings for a security policy
resource "akamai_appsec_advanced_settings_logging" "policy_logging" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  logging =  file("${path.module}/logging.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `logging` - (Required) The logging settings to apply ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#puthttpheaderloggingforaconfiguration)).

* `security_policy_id` - (Optional) The ID of a specific security policy to which the logging settings should be applied. If not supplied, the indicated settings will be applied to all policies within the configuration.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

