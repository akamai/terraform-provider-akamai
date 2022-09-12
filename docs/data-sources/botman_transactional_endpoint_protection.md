---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_transactional_endpoint_protection

**Scopes**: Security configuration

Returns information about the transactional endpoint protection settings assigned to a security configuration.

Use the [akamai_botman_transactional_endpoint_protection](../resources/akamai_botman_transactional_endpoint_protection) resource to create or modify your transactional endpoint protection settings.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/transactional-endpoint-protection](https://techdocs.akamai.com/bot-manager/reference/get-transactional-endpoint-protection). Returns information about all your transactional endpoints.

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

data "akamai_botman_transactional_endpoint_protection" "endpoint_protection" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "endpoint_protection_json" {
  value = data.akamai_botman_transaction_endpoint_protection.endpoint_protection.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the transactional endpoint protection settings.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your transactional endpoint protection settings.
