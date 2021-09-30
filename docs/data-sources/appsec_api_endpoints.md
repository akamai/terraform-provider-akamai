---
layout: "akamai"
page_title: "Akamai: ApiEndpoints"
subcategory: "Application Security"
description: |-
 ApiEndpoints
---

# akamai_appsec_api_endpoints

**Scopes**: Security configuration; security policy

Returns information about the API endpoints associated with a security policy or configuration. The returned information is described in the [Endpoint members](https://developer.akamai.com/api/cloud_security/application_security/v1.html#apiendpoint) section of the Application Security API documentation.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/api-endpoints](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapiendpoints)

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

data "akamai_appsec_api_endpoints" "api_endpoints" {
  config_id = 58843
  api_name  = "Contracts"
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the API endpoints.
- `security_policy_id` (Optional). Unique identifier of the security policy associated with the API endpoints. If not included, information is returned for all your security policies.
- `api_name` (Optional). Name of the API endpoint you want to return information for. If not included, information is returned for all your API endpoints.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `id_list`. List of API endpoint IDs.
- `json`. JSON-formatted list of information about the API endpoints.
- `output_text`. Tabular report showing the ID and name of the API endpoints.

