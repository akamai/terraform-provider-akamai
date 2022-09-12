---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_management_settings

**Scopes**: Security policy

Returns information about the bot management settings applied to a security policy.

To modify your existing bot management settings, use the [akamai_botman_bot_management_settings](../resources/akamai_botman_bot_management_settings) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/bot-management-settings](https://techdocs.akamai.com/bot-manager/reference/get-bot-management-settings). Returns your bot management settings and setting values.

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_botman_bot_management_settings" "management_settings" {
  config_id          =  data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "management_settings_json" {
  value = data.akamai_botman_bot_management_settings.management_settings.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot management settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot management settings.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your bot management settings and how they’re configured.
