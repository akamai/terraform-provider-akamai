---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_serve_alternate_action

**Scopes**: Security configuration; serve alternate action

Returns information about your serve alternate actions. 

Use the [akamai_botman_serve_alternate_action](../resources/akamai_botman_serve_alternate_action) resource to create or modify a serve alternate action.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/serve-alternate-actions](https://techdocs.akamai.com/bot-manager/reference/get-serve-alternate-actions). Returns information about all your serve alternate actions.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/serve-alternate-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/get-serve-alternate-action). Returns information about the specified serve alternate action.

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

// USE CASE: User wants to return information for all the serve alternate actions in the specified security configuration

data "akamai_botman_serve_alternate_action" "serve_alternate_actions" {
  config_id  = data.akamai_appsec_configuration.configuration.config_id
}

output "serve_alternate_actions_json" {
  value = data.akamai_botman_serve_alternate_action.serve_alternate_actions.json
}

// USE CASE: User only wants to return information for the serve alternate action with the ID cc9c3f89-e179-4892-89cf-d5e623ba9dc7

data "akamai_botman_serve_alternate_action" "serve_alternate_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  action_id = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}

output "serve_alternate_action_json" {
  value = data.akamai_botman_serve_alternate_action.serve_alternate_action.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the serve alternate actions.
- `action_id` (Optional). Unique identifier of the serve alternate action you want to return. If omitted, all your serve alternate actions are returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your serve alternate actions.

**See also**:

- [Set up alternate content](https://techdocs.akamai.com/bot-manager/docs/set-alternate-content)
