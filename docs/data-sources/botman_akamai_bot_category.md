---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_akamai_bot_category

**Scopes**: Universal (all bot categories defined by Akamai); bot category

Returns information about the bot categories [predefined by Akamai](https://techdocs.akamai.com/bot-manager/docs/akamai-categorized-bots). 

By including the `category_name` argument you can limit the returned data to a single category.

**Related API Endpoints**:

- [/appsec/v1/akamai-bot-categories](https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-categories). Returns information from all categories.
- [/appsec/v1/akamai-bot-categories/{categoryId}](https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-categories). Returns information only from the specified category.

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = “akamai/akamai”
    }
  }
}

provider “akamai” {
  edgerc = “~/.edgerc”
}

// USE CASE: User wants to return information for all bot categories

data “akamai_botman_akamai_bot_category” “bot_categories” {
}

output “bot_categories_json” {
  value = data.akamai_botman_akamai_bot_category.bot_categories.json
}

// USE CASE: User wants to return information for only the Web Search Engine Bots category

data “akamai_botman_akamai_bot_category” “bot_category” {
  category_name = “Akamai Bot Category 1”
}

output “bot_category_json” {
  value = data.akamai_botman_akamai_bot_category.bot_category.json
}
```

## Argument Reference

This resource supports the following arguments:

- `category_name` (Optional). Unique name of the Akamai bot category you want to return information for.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your Akamai-defined bot categories. The returned data includes the ID of each bot assigned to a category.

**See also**:

- [Akamai-categorized bots](https://techdocs.akamai.com/bot-manager/docs/akamai-categorized-bots)
