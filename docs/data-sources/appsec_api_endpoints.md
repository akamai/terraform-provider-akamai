---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_api_endpoints

**Scopes**: Security configuration; security policy

Returns information about the API endpoints associated with a security policy or configuration. 

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/api-endpoints](https://techdocs.akamai.com/application-security/reference/get-api-endpoints)

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