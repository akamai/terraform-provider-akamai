---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_defined_bot

**Scopes**: Security configuration; custom-defined bot

Creates or updates a custom-defined bot. To configure a custom bot you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `custom_defined_bot` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current set of custom-defined bots, use the [akamai_botman_custom_defined_bot](../data-sources/akamai_botman_custom_defined_bot) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-defined-bots](https://techdocs.akamai.com/bot-manager/reference/post-custom-defined-bot). Creates a new custom-defined bot.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-defined-bots/{botId}](https://techdocs.akamai.com/bot-manager/reference/put-custom-defined-bot). Updates an existing custom-defined bot.

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

// USE CASE: User wants to create a new custom bot

resource "akamai_botman_custom_defined_bot" "custom_defined_bot" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  custom_defined_bot = file("${path.module}/custom_defined_bot.json")
}

// USE CASE: User wants to modify the existing bot with the ID e08a628e-8a3c-4cd3-a5c9-8767c064ceb8

resource "akamai_botman_custom_defined_bot" "custom_defined_bot" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  custom_defined_bot = file("${path.module}/custom_defined_bot.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom defined bot.
- `custom_defined_bot` (Required). JSON collection of settings and setting values for the custom bot.  In the preceding sample code, the syntax `file("${path.module}/custom_defined_bot.json")` points to the location of a JSON file containing the custom bot settings and values.

**See also**:

- [Categorize and define your own bots](https://techdocs.akamai.com/bot-manager/docs/categorize-define-own-bots)
