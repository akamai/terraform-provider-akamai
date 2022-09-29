---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_conditional_action (Beta)

**Scopes**: Security configuration; conditional action

Returns information for your conditional actions. Conditional actions are actions typically designed to trigger in highly-specific situations.

Use the `action_id` argument to limit the returned data to information about the specified action.

To create or modify a conditional action, use the [akamai_botman_conditional_action](../resources/akamai_botman_conditional_action) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/conditional-actions](https://techdocs.akamai.com/bot-manager/reference/get-conditional-actions). Returns information about all your conditional actions.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/conditional-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/get-conditional-action). Returns information only for the specified conditional action.

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

// USE CASE: User wants to return information about all the conditional actions in the specified security configuration

data "akamai_botman_conditional_action" "conditional_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "conditional_actions_json" {
  value = data.akamai_botman_conditional_action.conditional_actions.json
}

// USE CASE: User only wants to return information for the conditional action with the ID cc9c3f89-e179-4892-89cf-d5e623ba9dc7

data "akamai_botman_conditional_action" "conditional_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  action_id = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}

output "conditional_action_json" {
  value = data.akamai_botman_conditional_action.conditional_action.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the conditional actions.
- `action_id` (Optional). Unique identifier of the conditional action you’d like returned. If omitted, all conditional actions are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your conditional actions and how they’re configured.

**See also**:

- [Set conditional actions](https://techdocs.akamai.com/bot-manager/docs/set-conditional-actions)
