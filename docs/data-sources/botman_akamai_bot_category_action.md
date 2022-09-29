---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_akamai_bot_category_action (Beta)

**Scopes**: Security policy; bot category

Returns information about the action taken when a bot category is triggered.

Use the `category_id` argument to return the action for a specified category. By default, information is returned for all Akamai-defined categories. And use the [akamai_botman_akamai_bot_category_action](../resources/akamai_botman_akamai_bot_category_action) resource to change the action assigned to an Akamai-defined bot category.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/akamai-bot-category-actions](https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-category-actions). Returns information for all Akamai bot categories.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/akamai-bot-category-actions/{categoryId}](https://techdocs.akamai.com/bot-manager/reference/get-akamai-bot-category-action). Returns information only for the specified Akamai bot category.

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

// USE CASE: User wants to return category actions for all bot categories

data "akamai_botman_akamai_bot_category_action" "category_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "category_actions_json" {
  value = data.akamai_botman_akamai_bot_category_action.category_actions.json
}

// USE CASE: User only wants to return information for the bot category with the ID 2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8

data "akamai_botman_akamai_bot_category_action" "category_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  category_id        = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
}

output "category_action_json" {
  value = data.akamai_botman_akamai_bot_category_action.category_action.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot category.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot category.
- `category_id` (Optional). Unique identifier of an Akamai-defined bot category. Use this argument if you want to return the action for a specific category.

## Output Options

The following options determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your bot categories and their assigned action. Each category can have only one assigned action, an action that applies to each bot in the category.

**See also**:

- [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot)
