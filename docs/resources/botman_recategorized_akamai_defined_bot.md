---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_recategorized_akamai_defined_bot (Beta)

**Scopes**: Security configuration

Moves an Akamai-defined bot to a custom bot category.

To review your current set of recategorized bots, use the [akamai_botman_recategorized_akamai_defined_bot](../data-sources/akamai_botman_recategorized_akamai_defined_bot) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/recategorized-akamai-defined-bots](https://techdocs.akamai.com/bot-manager/reference/post-recategorized-akamai-defined-bot). Recategorizes the specified Akamai-defined bot.

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

resource "akamai_botman_recategorized_akamai_defined_bot" "recategorized_bot" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  bot_id      = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  category_id = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot being recategorized.
- `bot_id` (Required). Unique identifier of the Akamai-defined bot to be recategorized.
- `category_id` (Required). Unique identifier of the custom category the bot is being moved to. Note that you can only move bots to a custom category. You can’t move a bot to an Akamai-defined category.
