---
layout: "akamai"
page_title: "Akamai: ApiEndpoints"
subcategory: "Application Security"
description: |-
 ApiEndpoints
---

# akamai_appsec_api_endpoints

Use the `akamai_appsec_api_endpoints` data source to retrieve information about the API Endpoints associated with a security policy or configuration version. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapiendpoints).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_api_endpoints" "api_endpoints" {
  config_id = 43253
  version = 7
  api_name = "TestEndpoint"
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID.

* `version` - (Required) The version number of the configuration.

* `security_policy_id` - (Optional) The ID of the security policy to use.

* `api_name` - (Optional) The name of a specific endpoint.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id_list` - A list of IDs of the API endpoints.

* `json` - A JSON-formatted list of information about the API endpoints.

* `output_text` - A tabular display showing the ID and name of the API endpoints.

