---
layout: "akamai"
page_title: "Akamai: ApiRequestConstraints"
subcategory: "Application Security"
description: |-
 ApiRequestConstraints
---

# akamai_appsec_api_request_constraints

Use the `akamai_appsec_api_request_constraints` data source to retrieve a list of APIs with their constraints and associated actions, or the constraints and actions for a particular API. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapirequestconstraints).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the all api request constraints associated with a given security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_api_request_constraints" "apis_request_constraints" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

//endpoint id and action
output "apis_constraints_text" {
  value = data.akamai_appsec_api_request_constraints.apis_request_constraints.output_text
}

output "apis_constraints_json" {
  value = data.akamai_appsec_api_request_constraints.apis_request_constraints.json
}

// USE CASE: user wants to view action on a single api request constraint associated with a given security policy
data "akamai_appsec_api_request_constraints" "api_request_constraints" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  api_id = var.api_id
}

output "api_constraints_text" {
  value = data.akamai_appsec_api_request_constraints.api_request_constraints.output_text
}

output "api_constraints_json" {
  value = data.akamai_appsec_api_request_constraints.api_request_constraints.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `api_id` - (Optional) The ID of a specific API for which to retrieve constraint information.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of information about the APIs and their constraints and actions.

* `output_text` - A tabular display showing the APIs and their constraints and actions.

