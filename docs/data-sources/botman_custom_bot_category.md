---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_bot_category (Beta)

**Scopes**: Security configuration; custom bot category

Returns information about the custom bot categories you’ve created.

By including the `category_id` argument you can limit the returned data to a single category.

To create or update a custom bot category, use the [akamai_botman_custom_bot_category](../resources/akamai_botman_custom_bot_category) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-bot-categories](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-categories). Returns information about all your custom categories.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-bot-categories/{categoryId}](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category). Returns the category action assigned to the specified custom category.

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

// USE CASE: User wants to return information for all custom categories in the specified security configuration

data "akamai_botman_custom_bot_category" "custom_categories" {
  config_id  = data.akamai_appsec_configuration.configuration.config_id
}

output "custom_categories_json" {
  value = data.akamai_botman_custom_bot_category.custom_category.json
}

// USE CASE: User only wants to return information for the custom category with the ID 2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8

data "akamai_botman_custom_bot_category" "custom_category" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  category_id = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
}

output "custom_category_json" {
  value = data.akamai_botman_custom_bot_category.custom_category.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom bot categories.
- `category_id` (Optional). Unique identifier of the custom bot category you want returned. If omitted, information about all your custom categories is returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your custom bot categories. This output includes the IDs of both the custom bots and the Akamai-defined bots assigned to a category.

**See also**:

- [Categorize and define your own bots](https://techdocs.akamai.com/bot-manager/docs/categorize-define-own-bots)
