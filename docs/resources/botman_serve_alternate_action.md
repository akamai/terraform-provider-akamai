---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_serve_alternate_action

**Scopes**: Security configuration; serve alternate action

Creates or updates a serve alternate action. 

To configure a serve alternate action you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `serve_alternate_action` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your existing serve alternate actions, use the [akamai_botman_serve_alternate_action](../data-sources/akamai_botman_serve_alternate_action) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/serve-alternate-actions](https://techdocs.akamai.com/bot-manager/reference/post-serve-alternate-action). Creates a serve alternate action.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/serve-alternate-actions/{actionId}](https://techdocs.akamai.com/bot-manager/reference/put-serve-alternate-action). Updates an existing serve alternate action.

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

// USE CASE: User wants to create a new serve alternate action

resource "akamai_botman_serve_alternate_action" "serve_alternate_action" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  serve_alternate_action = file("${path.module}/serve_alternate_action.json")
}

//USE CASE: User wants to modify the serve alternate action with the ID akamai

resource "akamai_botman_serve_alternate_action" "serve_alternate_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  serve_alternate_action = file("${path.module}/serve_alternate_action.json")
EOF
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the serve alternate action.
- `serve_alternate_action` (Required). JSON-formatted collection of serve alternate action settings and setting values. In the preceding sample code, the syntax `file("${path.module}/serve_alternate_action.json")` points to the location of a JSON file containing the serve alternate action settings and values.

**See also**:

- [Set up alternate content](https://techdocs.akamai.com/bot-manager/docs/set-alternate-content#:~:text=Configurable%20actions-,Set%20up%20alternate%20content,-Set%20conditional%20actions)
