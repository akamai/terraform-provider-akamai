---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsLogging"
subcategory: "Application Security"
description: |-
  AdvancedSettingsLogging
---

# akamai_appsec_advanced_settings_logging

**Scopes**: Security configuration; security policy

Enables, disables, or updates HTTP header logging settings.
By default, this operation applies at the configuration level, which means that it applies to all the security policies within that configuration.
However, by using the `security_policy_id` parameter you can specify custom settings for an individual security policy.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/logging](https://techdocs.akamai.com/application-security/reference/put-policies-logging)

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

// USE CASE: User wants to modify the logging settings for a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_advanced_settings_logging" "logging" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  logging   = file("${path.module}/logging.json")
}

// USE CASE: User wants to configure logging settings for a security policy.

resource "akamai_appsec_advanced_settings_logging" "policy_logging" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  logging            = file("${path.module}/logging.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration containing the logging settings being modified.
- `logging` (Required). Path to a JSON file containing the logging settings to be configured. A sample JSON file can be found in the [Modify HTTP header log settings for a configuration](https://developer.akamai.com/api/cloud_security/application_security/v1.html#puthttpheaderloggingforaconfiguration) section of the Application Security API documentation.
- `security_policy_id` (Optional). Unique identifier of the security policies whose settings are being modified. If not included, the logging settings are modified at the configuration scope and, as a result, apply to all the security policies associated with the configuration.