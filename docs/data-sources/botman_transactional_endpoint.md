---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_transactional_endpoint

**Scopes**: Security policy; operation

Returns information about your transactional endpoints. 

To create or modify a transactional endpoint, use the [akamai_botman_transactional_endpoint](../resources/akamai_botman_transactional_endpoint) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/transactional-endpoints/bot-protection](https://techdocs.akamai.com/bot-manager/reference/get-transactional-endpoints). Returns information for all your transactional endpoints.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/transactional-endpoints/bot-protection/{operationId}](https://techdocs.akamai.com/bot-manager/reference/get-transactional-endpoint). Returns information for the specified transactional endpoint.

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

// USE CASE: User wants to return information for all the transactional endpoints in the specified security configuration

data "akamai_botman_transactional_endpoint" "transactional_endpoint" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "transactional_endpoint_json" {
  value = data.akamai_botman_transactional_endpoint.transactional_endpoint.json
}

// USE CASE: User only wants to return data for the transactional endpoint with the ID e0f89bb0-77d5-46f7-979d-e204e6fdc5a5

data "akamai_botman_transactional_endpoint" "transactional_endpoint" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  operation_id       = "e0f89bb0-77d5-46f7-979d-e204e6fdc5a5"
}

output "transactional_endpoint_json" {
  value = data.akamai_botman_transactional_endpoint.transactional_endpoint.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the transactional endpoints.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the transactional endpoints.
- `operation_id` (Optional). Unique identifier of the API operation to be returned. If omitted, transactional endpoint information is returned for all your operations.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your transactional endpoints.
