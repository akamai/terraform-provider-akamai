---
layout: "akamai"
page_title: "Akamai: ApiRequestConstraints"
subcategory: "Application Security"
description: |-
 ApiRequestConstraints
---

# akamai_appsec_api_request_constraints

**Scopes**: Security policy; API endpoint

Returns information about API endpoint constraints and actions. 

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/api-request-constraints](https://techdocs.akamai.com/application-security/reference/get-api-request-constraints)

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

// USE CASE: User wants to view all the API request constraints associated with a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_api_request_constraints" "apis_request_constraints" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "apis_constraints_text" {
  value = data.akamai_appsec_api_request_constraints.apis_request_constraints.output_text
}

output "apis_constraints_json" {
  value = data.akamai_appsec_api_request_constraints.apis_request_constraints.json
}

// USE CASE: User wants to view the action associated with an API request constraint.

data "akamai_appsec_api_request_constraints" "api_request_constraints" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  api_id             = 624913
}

output "api_constraints_text" {
  value = data.akamai_appsec_api_request_constraints.api_request_constraints.output_text
}

output "api_constraints_json" {
  value = data.akamai_appsec_api_request_constraints.api_request_constraints.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the API constraints.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the API constraints.
- `api_id` (Optional). Unique identifier of the API endpoint you want to return constraint information for.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of information about the APIs, their constraints, and their actions.
- `output_text`. Tabular report of the APIs, their constraints, and their actions.