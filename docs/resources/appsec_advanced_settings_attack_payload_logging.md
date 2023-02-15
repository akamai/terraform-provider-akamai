---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_advanced_settings_attack_payload_logging

**Scopes**: Security configuration; security policy

Enables, disables, or updates Attack Payload Logging settings.
By default, this operation is applied at the configuration level, which means that it is applied to all the security policies within that configuration.
However, by using the `security_policy_id` parameter you can specify custom settings for an individual security policy.

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

// USE CASE: User wants to modify the Attack Payload Logging settings for a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_advanced_settings_attack_payload_logging" "attack_payload_logging" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  attack_payload_logging   = file("${path.module}/attack-payload-logging.json")
}

// USE CASE: User wants to configure Attack Payload Logging settings for a security policy.

resource "akamai_appsec_advanced_settings_attack_payload_logging" "policy_logging" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  attack_payload_logging            = file("${path.module}/attack-payload-logging.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration containing the Attack Payload Logging settings being modified.
- `attack_payload_logging` (Required). JSON representation of the Attack Payload Logging settings to be configured.
- `security_policy_id` (Optional). Unique identifier of the security policies whose settings are being modified. If not included, the Attack Payload Logging settings are modified at the configuration scope and, as a result, apply to all the security policies associated with the configuration.