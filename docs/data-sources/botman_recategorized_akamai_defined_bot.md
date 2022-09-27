---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_recategorized_akamai_defined_bot (Beta)

**Scopes**: Security configuration; bot

Returns information about your recategorized bots. A recategorized bot is an Akamai-defined bot that’s been moved from an Akamai-defined category to a custom category.

Use the [akamai_botman_recategorized_akamai_defined_bot](../resources/akamai_botman_recategorized_akamai_defined_bot) resource to recategorize a bot.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/recategorized-akamai-defined-bots](https://techdocs.akamai.com/bot-manager/reference/get-recategorized-akamai-defined-bots). Returns information about all your recategorized bots.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/recategorized-akamai-defined-bots/{botId}](https://techdocs.akamai.com/bot-manager/reference/get-recategorized-akamai-defined-bot). Returns information about the specified recategorized bot.

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

// USE CASE: User wants to return information for all the recategorized bots in the specified security configuration

data "akamai_botman_recategorized_akamai_defined_bot" "recategorized_bots" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "recategorized_bots_json" {
  value = data.akamai_botman_recategorized_akamai_defined_bot.recategorized_bot.json
}

// USE CASE: User only wants to return information for the recategorized bot with the ID cc9c3f89-e179-4892-89cf-d5e623ba9dc7

data "akamai_botman_recategorized_akamai_defined_bot" "recategorized_bot" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  bot_id    = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}

output "recategorized_bot_json" {
  value = data.akamai_botman_recategorized_akamai_defined_bot.recategorized_bot.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the recategorized bot.
- `bot_id` (Optional). Unique identifier of the recategorized bot you want returned. If omitted, all your recategorized bots are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about the Akamai-defined bots that have been moved to a custom bot category.
