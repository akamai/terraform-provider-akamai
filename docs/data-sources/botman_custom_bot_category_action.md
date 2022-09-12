---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_bot_category_action

**Scopes**: Security policy; custom category action

Returns information about the action taken when a custom bot category is triggered. 

Use the `category_id` argument to return the category action for a specified category. By default, information is returned for all the category actions associated with a specific security policy.

Use the [akamai_botman_custom_bot_category_action](../resources/akamai_botman_custom_bot_category_action) resource to modify the action assigned to a custom bot category.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/custom-bot-category-actions](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-actions). Returns all your custom categories and the actions assigned to them.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/custom-bot-category-actions/{categoryId}](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-action). Returns the action assigned to the specified custom category.

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

// USE CASE: User wants to return information for all the custom category actions

data "akamai_botman_custom_bot_category_action" "category_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "category_actions_json" {
  value = data.akamai_botman_custom_bot_category_action.category_actions.json
}

// USE CASE: User only wants to return category actions for the custom category with the ID 2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8

data "akamai_botman_custom_bot_category_action" "category_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  category_id        = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
}

output "category_actions_json" {
  value = data.akamai_botman_custom_bot_category_action.category_actions.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom bot category.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the custom bot category.
- `category_id` (Optional). Unique identifier of the custom category whose action is being returned. If omitted, actions for all your custom categories are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your custom categories and the action assigned to them.

**See also**:

- [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot)
