---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_defined_bot (Beta)

**Scopes**: Security configuration; custom bot

Returns information about your custom-defined bots.

Use the `bot_id` argument to limit the returned data to the specified bot.

To create or modify a custom-defined bot, use the [akamai_botman_custom_defined_bot](../resources/akamai_botman_custom_defined_bot) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-defined-bots](https://techdocs.akamai.com/bot-manager/reference/get-custom-defined-bots). Returns information about all your custom-defined bots.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-defined-bots/{botId}](https://techdocs.akamai.com/bot-manager/reference/get-custom-defined-bot). Returns information about the specified custom-defined bots.

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

// USE CASE: User wants to return information for all custom-defined bots

data "akamai_botman_custom_defined_bot" "custom_bots" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "custom_bots_json" {
  value = data.akamai_botman_custom_defined_bot.custom_bots.json
}

// USE CASE: User only wants to return data for the custom-defined bot with the ID e08a628e-87dc-4343-a5c9-8767c061ceb8

data "akamai_botman_custom_defined_bot" "custom_bot" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  bot_id    = "e08a628e-87dc-4343-a5c9-8767c061ceb8"
}

output "custom_bots_json" {
  value = data.akamai_botman_custom_defined_bot.custom_bot.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom-defined bots.
- `bot_id` (Optional). Unique identifier of the custom bot you want returned. If omitted, all your custom bots are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your custom bots.

**See also**:

- [Categorize and define your own bots](https://techdocs.akamai.com/bot-manager/docs/categorize-define-own-bots)
