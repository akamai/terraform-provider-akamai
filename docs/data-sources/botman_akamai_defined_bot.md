---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_akamai_defined_bot (Beta)

**Scopes**: Universal (all bots defined by Akamai); bot

Returns information about the bots predefined by Akamai. Use the `bot_name` argument to limit the returned data to a specific bot.

**Related API Endpoints**:

- [/appsec/v1/akamai-defined-bots](https://techdocs.akamai.com/bot-manager/reference/get-akamai-defined-bots-1). Returns data for all Akamai-defined bots.
- [/appsec/v1/akamai-defined-bots/{botId}](https://techdocs.akamai.com/bot-manager/reference/get-akamai-defined-bot-1). Returns data for the specified Akamai-defined bot.

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

// USE CASE: User wants to return information for all Akamai-defined bots

data "akamai_botman_akamai_defined_bot" "defined_bots" {
}
output "defined_bots_json" {
  value = data.akamai_botman_akamai_defined_bot.defined_bots.json
}

// USE CASE: User only wants to return information for the price-scraper-bot bot

data "akamai_botman_akamai_defined_bot" "defined_bot" {
  bot_name = "price-scraper-bot"
}

output " data "defined_bot_json" {
  value = data.akamai_botman_akamai_defined_bot.defined_bot.json
}
```

## Argument Reference

This resource supports the following arguments:

- `bot_name` (Optional). Unique name of an Akamai-defined bot.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

-`json`. JSON-formatted output containing information about one or more Akamai-defined bots. The returned information includes the bot name, bot ID, and category ID.

**See also**:

- [Akamai-categorized bots](https://techdocs.akamai.com/bot-manager/docs/akamai-categorized-bots)
