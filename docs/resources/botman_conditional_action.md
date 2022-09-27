---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_conditional_action (Beta)

**Scopes**: Security configuration; conditional action

Creates or updates a conditional action.

To configure a conditional action you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `conditional_action` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your existing conditional actions, use the [akamai_botman_conditional_action](../data-sources/akamai_botman_conditional_action) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/conditional-actions](https://techdocs.akamai.com/bot-manager/reference/post-conditional-action). Creates a new conditional action.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/conditional-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/put-conditional-action). Updates an existing conditional action.

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

// USE CASE: User wants to create a new conditional action

resource "akamai_botman_conditional_action" "conditional_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  conditional_action = file("${path.module}/conditional_action.json")
}

// USE CASE: User wants to update an existing conditional action

resource "akamai_botman_conditional_action" "conditional_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  conditional_action = file("${path.module}/conditional_action.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the conditional actions.
- `conditional_action` (Required). JSON-formatted collection of conditional action settings and setting values. In the preceding sample code, the syntax `file("${path.module}/conditional_action.json")` points to the location of a JSON file containing the conditional action settings and values.

••See also**:

- [Set up conditional actions](https://techdocs.akamai.com/bot-manager/docs/set-conditional-actions)
