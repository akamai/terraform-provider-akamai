---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_challenge_action (Beta)

**Scopes**: Security configuration; challenge action

Creates or updates a challenge action.

To configure a challenge action you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `challenge_action` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current challenge actions use the [akamai_botman_challenge_action](../data-sources/akamai_botman_challenge_action) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/challenge-actions](https://techdocs.akamai.com/bot-manager/reference/post-challenge-action-1). Creates a new challenge action.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/challenge-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/put-challenge-action-1). Updates the specified challenge action.

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

// USE CASE: User wants to create a new challenge action challenge_108949

resource "akamai_botman_challenge_action" "challenge_action" {
  config_id        = data.akamai_appsec_configuration.configuration.config_id
  challenge_action = file("${path.module}/challenge_action.json")
}

// USE CASE: Users wants to update the challenge action with the ID challenge_108949

resource "akamai_botman_challenge_action" "challenge_action" {
  config_id        = data.akamai_appsec_configuration.configuration.config_id
  challenge_action = file("${path.module}/challenge_action.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the challenge action.
- `challenge_action` (Required). JSON-formatted collection of challenge action settings and their values. In the preceding sample code, the syntax `file("${path.module}/challenge_action.json")` points to the location of a JSON file containing the challenge action settings and values.

**See also**:

- [Challenge actions](https://techdocs.akamai.com/bot-manager/docs/challenge-actions)
