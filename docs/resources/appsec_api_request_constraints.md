---
layout: "akamai"
page_title: "Akamai: ApiRequestConstraints"
subcategory: "Application Security"
description: |-
  ApiRequestConstraints
---

# resource_akamai_appsec_api_request_constraints

The `resource_akamai_appsec_api_request_constraints` resource allows you to update what action to take when the API request constraint triggers. This operation modifies an individual API constraint action. To use this operation, use the `akamai_appsec_api_endpoints` data source to list one or all API endpoints, and use the ID of the selected endpoint. Use the `action` paameter to specify how the alert should be handled.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the api request constraints action
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_api_endpoints" "api_endpoint" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  name = var.api_endpoint_name
}

resource "akamai_api_request_constraints" "api_request_constraints" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  api_endpoint_id = data.akamai_appsec_api_endpoints.api_endpoint.id
  action = "alert"
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `api_endpoint_id` - (Required) The ID of the API endpoint to use.

* `action` - (Required) The action to assign to API request constraints: either `alert`, `deny`, or `none`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

