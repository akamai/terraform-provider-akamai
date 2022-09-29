---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_response_action (Beta)

**Scopes**: Security configuration; response action

Returns information about the actions that can be taken when a bot detection method is triggered.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions](https://techdocs.akamai.com/bot-manager/reference/get-response-actions). Returns information about all your response actions.

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

// USE CASE: User wants to return information about all the response actions in the specified security configuration

data "akamai_botman_response_action" "response_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "response_actions_json" {
  value = data.akamai_botman_response_action.response_actions.json
}

// USE CASE: User only wants to return information about the response action with the ID cc9c3f89-e179-4892-89cf-d5e623ba9dc7

data "akamai_botman_response_action" " response_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  action_id = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}

output "response_action_json" {
  value = data.akamai_botman_response_action.response_action.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the response actions.
- `action_id` (Optional). Unique identifier of the response action you want returned. If omitted, all your response actions are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your response actions.

**See also**:

- [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot)
