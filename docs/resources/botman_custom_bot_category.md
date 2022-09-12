---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_bot_category

**Scopes**: Security configuration; custom bot category

Creates or modifies a bot category that you can use in addition to the Akamai-defined bot categories. 

To configure a custom category you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `custom category` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current set of custom bot categories, use the [akamai_botman_custom_bot_category](../data-sources/akamai_botman_custom_bot_category) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-bot-categorieS](https://techdocs.akamai.com/bot-manager/reference/post-custom-bot-category). Creates a custom bot category.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-bot-categories/{categoryId}](https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category). Updates an existing custom bot category.

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

// USE CASE: User wants to create a new custom bot category

resource "akamai_botman_custom_bot_category" "custom_bot_category" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  custom_bot_category = file("${path.module}/custom_bot_category.json")
}

// USE CASE: User wants to modify the custom bot category with the ID a08a6d8e-a23c-4cb3-a5c9-8e3dc0d4c0b8

resource "akamai_botman_custom_bot_category" "custom_bot_category" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  custom_bot_category = file("${path.module}/custom_bot_category.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom bot category.
- `custom_bot_category` (Required). JSON-formatted collection of bot settings and setting values.  In the preceding sample code, the syntax `file("${path.module}/custom_bot_category.json")` points to the location of a JSON file containing the custom category settings and values.

**See also**:

- [Categorize and define your own bots](https://techdocs.akamai.com/bot-manager/docs/categorize-define-own-bots)
