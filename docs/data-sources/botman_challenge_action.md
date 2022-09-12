---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_challenge_action

Scopes: Security configuration; challenge action

Returns information about your bot challenge actions. A challenge action is a process (such as ReCAPTCHA) that must be completed successfully before a request can be processed.

Use the `action_id` argument to limit the returned data to information about a specific challenge action.

To create or modify a challenge action, use the [akamai_botman_challenge_action](../resources/akamai_botman_challenge_action) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/challenge-actions](https://techdocs.akamai.com/bot-manager/reference/get-challenge-actions-1). Returns information about all your challenge actions.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/challenge-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/get-challenge-action-1). Returns information about the specified challenge action.

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

// USE CASE: User wants to return all challenge actions for the specified security configuration

data "akamai_botman_challenge_action" "challenge_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "challenge_actions_json" {
  value = data.akamai_botman_challenge_action.challenge_actions.json
}

// USE CASE: User only wants to return information for challenge action cc9c3f89-e179-4892-89cf-d5e623ba9dc7 in the specified security configuration

data "akamai_botman_challenge_action" "challenge_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  action_id = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}

output "challenge_action_json" {
  value = data.akamai_botman_challenge_action.challenge_action.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the challenge actions.
- `action_id` (Optional). Unique identifier of the challenge action to be returned. If omitted, information is returned for all your challenge actions.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your challenge actions and how they’re configured.

**See also**:

- [Challenge actions](https://techdocs.akamai.com/bot-manager/docs/challenge-actions)
