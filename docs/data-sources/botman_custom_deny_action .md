---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_deny_action

**Scopes**: Security configuration; custom deny action

Returns information about your custom deny actions. 

> **Note**. Custom deny actions aren’t available for Akamai’s ChinaCDN.

To create or modify a custom deny action, use the [akamai_botman_custom_deny_action](../resources/akamai_botman_custom_deny_action) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/custom-deny-actions](https://techdocs.akamai.com/bot-manager/reference/get-custom-deny-actions). Returns information about all your custom deny actions.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/custom-deny-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/get-custom-deny-action). Returns information about the specified custom deny action.

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

// USE CASE: User wants to return information about all custom deny actions in the specified security configuration

data "akamai_botman_custom_deny_action" "custom_deny_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "custom_deny_actions_json" {
  value = data.akamai_botman_custom_deny_action.custom_deny_actions.json
}

// USE CASE: User only wants to return information for the custom deny action with the ID cc9c3f89-e179-4892-89cf-d5e623ba9dc7

data "akamai_botman_custom_deny_action" "custom_deny_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  action_id = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}

output "custom_deny_actions_json" {
  value = data.akamai_botman_custom_deny_action.custom_deny_actions.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom deny actions.
- `action_id` (Optional). Unique identifier of the custom deny action you want returned. If omitted, all your custom deny actions are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your custom deny actions.

**See also**:

- [Create a custom deny action] https://techdocs.akamai.com/bot-manager/docs/create-custom-deny-action)
