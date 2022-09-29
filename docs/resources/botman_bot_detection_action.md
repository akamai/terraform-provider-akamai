---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_detection_action (Beta)

**Scopes**: Bot detection method

Modifies the action assigned to a bot detection method.

To review your current bot detection actions, use the [akamai_botman_bot_detection_action](../data-sources/akamai_botman_bot_detection_action) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/bot-detection-actions/{detectionId}](https://techdocs.akamai.com/bot-manager/reference/put-bot-detection-action). Updates the specified bot detection action.

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

resource "akamai_botman_bot_detection_action" "bot_detection_action" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id   = "gms1_134637"
  detection_id         = "65cd1c42-d8e9-42af-9f78-153cfdd92443"
  bot_detection_action = file("${path.module}/action.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot detection method.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot detection method.
- `detection_id` (Required). Unique identifier of the bot detection method being updated.
- `bot_detection_action` (Required). JSON file containing the action taken when the bot detection method is triggered.

**See also**:

- [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot)
