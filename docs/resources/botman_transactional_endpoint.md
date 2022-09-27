---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_transactional_endpoint (Beta)

**Scopes**: Security policy; API operation

Creates or updates a transactional endpoint. To configure a transactional endpoint you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `transactional_endpoint` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your existing transactional endpoints, use the [akamai_botman_transactional_endpoint](../data-sources/akamai_botman_transactional_endpoint) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/transactional-endpoints/bot-protection](https://techdocs.akamai.com/bot-manager/reference/post-transactional-endpoint). Creates a transactional endpoint.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/transactional-endpoints/bot-protection/{operationId}](https://techdocs.akamai.com/bot-manager/reference/put-transactional-endpoint). Updates a transactional endpoint.

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

// USE CASE: User wants to create a new transactional endpoint

resource "akamai_botman_transactional_endpoint" "transaction_endpoint" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = "gms1_134637"
  transactional_endpoint = file("${path.module}/transactional_endpoint.json")
  operation_id           = "e0f89bb0-77d5-46f7-979d-e204e6fdc5a5"
}

// USE CASE: User wants to update the transactional endpoint with the operation ID adbe9bb8-732f-4935-a725-09dd2dbe66dd

resource "akamai_botman_transactional_endpoint" "transactional_endpoint" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = "gms1_134637"
  operation_id           = "e0f89bb0-77d5-46f7-979d-e204e6fdc5a5"
  transactional_endpoint = file("${path.module}/transactional_endpoint.json")
EOF
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the transactional endpoint.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the transactional endpoint.
- `operation_id` (Required). Unique identifier of the API operation being created or updated.
- `transactional_endpoint` (Required). JSON collection of transactional endpoint settings and setting values. In the preceding sample code, the syntax `file("${path.module}/transactional_endpoint.json")` points to the location of a JSON file containing the transactional endpoint settings and values.
