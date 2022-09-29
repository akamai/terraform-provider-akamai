---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_deny_action (Beta)

**Scopes**: Security configuration; custom deny action

Creates or modifies a custom deny action.

To configure a custom deny action you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `custom_deny_action` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current set of custom deny actions, use the [akamai_botman_custom_deny_action](../data-sources/akamai_botman_custom_deny_action) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/custom-deny-actions](https://techdocs.akamai.com/bot-manager/reference/post-custom-deny-action). Creates a custom deny action.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/custom-deny-actions/{actionId}] https://techdocs.akamai.com/bot-manager/reference/put-custom-deny-action). Updates an existing custom deny action.

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

// USE CASE: User wants to create a new custom deny action

resource "akamai_botman_custom_deny_action" "custom_deny_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  custom_deny_action = file("${path.module}/custom_deny_action.json")
}

//USE CASE: User wants to update an existing custom deny action

resource "akamai_botman_custom_deny_action" "custom_deny_action" {
  config_id          = 43253
  custom_deny_action = file("${path.module}/custom_deny_action.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom deny action.
- `custom_deny_action` (Required). JSON-formatted collection of custom deny action settings and setting values. In the preceding sample code, the syntax `file("${path.module}/custom_deny_action.json")` points to the location of a JSON file containing the custom deny action settings and values.

**See also**:

- [Create a custom deny action](https://techdocs.akamai.com/bot-manager/docs/create-custom-deny-action)
