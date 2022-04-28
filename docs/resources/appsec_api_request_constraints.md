---
layout: "akamai"
page_title: "Akamai: ApiRequestConstraints"
subcategory: "Application Security"
description: |-
  ApiRequestConstraints
---

# akamai_appsec_api_request_constraints

**Scopes**: API endpoint

Modifies the action taken when an API request constraint triggers.
To use this operation, call the [akamai_appsec_api_endpoints](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_attack_group) data source to list the names of your API endpoints, then apply a constraint (and an accompanying action) to one of those endpoints.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/api-request-constraints/{apiId}](https://techdocs.akamai.com/application-security/reference/put-api-request-constraints-api)

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

// USE CASE: User wants to set the API request constraints action.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_api_endpoints" "api_endpoint" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  api_name           = "Contracts"
}

resource "akamai_appsec_api_request_constraints" "api_request_constraints" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  api_endpoint_id    = data.akamai_appsec_api_endpoints.api_endpoint.id
  action             = "alert"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the API request constraint settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the API request constraint settings being modified.
- `api_endpoint_id` (Optional). ID of the API endpoint the constraint will be assigned to.
- `action` (Required). Action to assign to the API request constraint. Allowed values are:
  - **alert**, Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.