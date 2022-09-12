---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_management_settings

**Scopes**: Security policy

Modifies the bot management settings for the specified security policy. To configure a Bot Manager settings you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `bot_management_settings` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current bot management settings, use the [akamai_botman_bot_management_settings](./data-sources/akamai_botman_bot_management_settings) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/bot-management-settings](https://techdocs.akamai.com/bot-manager/reference/put-bot-management-settings). Updates the bot management settings for the specified security policy.

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

resource "akamai_botman_bot_management_settings" "bot_management_settings" {
  config_id               = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id      = "gms1_134637"
  bot_management_settings = file("${path.module}/bot_management_settings.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot management settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot management settings.
- `bot_management_settings` (Required). JSON-formatted collection of bot management settings and their values. In the preceding sample code, the syntax `file("${path.module}/bot_management_settings.json")` points to the location of a JSON file containing the Bot Manager settings and values.
