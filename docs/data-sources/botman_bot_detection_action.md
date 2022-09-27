---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_detection_action (Beta)

**Scopes**: Security policy; bot detection action

Returns information about the action taken when a bot detection method is triggered.

Use the `detection_id` argument to return the action for a specified detection method. By default, information is returned for all the detection methods associated with a specific security policy. And use the [akamai_botman_bot_detection_action](../resources/akamai_botman_bot_detection_action) resource to modify an existing bot detection action.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/bot-detection-actions]( https://techdocs.akamai.com/bot-manager/docs/bot-det-methods-rule-ids). Returns action information for all your detection actions.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/bot-detection-actions/{detectionId}](https://techdocs.akamai.com/bot-manager/reference/get-bot-detection-action). Returns detection action information for the specified detection method.

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

// USE CASE: User wants to return all detection actions for the specified security policy

data "akamai_botman_bot_detection_action" "detection_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "detections_actions_json" {
  value = data.akamai_botman_bot_detection_action.detection_actions.json
}

// USE CASE: User only wants to return information for the detection action with the ID 65cd1c42-d8e9-42af-9f78-153cfdd92443

data "akamai_botman_bot_detection_action" "detection_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  detection_id       = "65cd1c42-d8e9-42af-9f78-153cfdd92443"
}

output "detection_action_json" {
  value = data.akamai_botman_bot_detection_action.detection_action.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot detection methods.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot detection methods.
- `detection_id` (Optional). Unique identifier of the detection method whose action is being returned. If omitted, information is returned for all the actions of all the detection methods associated with the security policy.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about the actions associated with your bot detection methods.

**See also**:

- [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot)
